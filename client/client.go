package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	pbf "gitlab.noncepad.com/naomiyoko/ffmpeg-market/proto/ffmpeg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	ProcessRequest(ctx context.Context, blender string, outDir string, extenstionList []string) error
}

type external struct {
	client pbf.JobManagerClient
}

// host:port
// args: 0=url string, 1=fileIn string, 2=DirOut string, 3=ext []string
func Run(parentCtx context.Context, args []string) error {

	go func() {
		<-parentCtx.Done()
		log.Print("parent ctx done")
	}()
	// Check if the output directory exists
	if _, err := os.Stat(args[2]); os.IsNotExist(err) {
		return fmt.Errorf("Run - 1: output directory does not exist, err: %s", err)
	}

	// Create a client
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	conn, err := grpc.DialContext(ctx, args[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client := Create(ctx, conn)
	blenderIn := args[1]
	dir := args[2]
	extList := args[3:]

	err = client.ProcessRequest(ctx, blenderIn, dir, extList)
	if err != nil {
		return err
	}

	return err
}

func Create(ctx context.Context, conn *grpc.ClientConn) Client {
	// Create a new gRPC client
	return external{client: pbf.NewJobManagerClient(conn)}
}

func outfilePath(dir string, ext string) string {
	return fmt.Sprintf("%s/out.%s", dir, ext)
}

func (e external) ProcessRequest(parentCtx context.Context, blender string, outDir string, extenstionList []string) error {
	// Check if the file exists and has the correct extension
	if _, err := os.Stat(blender); os.IsNotExist(err) {
		return fmt.Errorf("processRequest - 1: file does not exist, err: %s", err)
	} else if filepath.Ext(blender) != ".blend" {
		return fmt.Errorf("process Request - 2: file must have .blend extension, err %s", err)
	}

	// check if file extenstions are compatible with ffmpeg
	var supportedExtensions = map[string]bool{
		"mp4":  true,
		"mov":  true,
		"mkv":  true,
		"flv":  true,
		"wmv":  true,
		"webm": true,
		"mpeg": true,
		"ogv":  true,
		"gif":  true,
	}
	for _, ext := range extenstionList {
		if !supportedExtensions[ext] {
			return fmt.Errorf("processRequest - 3: extension %s is not compatible with ffmpeg", ext)
		}
	}
	// set context with cancel
	ctx, cancel := context.WithCancelCause(parentCtx)

	// gRPC call to the server
	stream, err := e.client.Process(ctx)
	if err != nil {
		cancel(err)
		return err
	}
	// create_args() w extensionlist
	args := createArgs(extenstionList)
	err = stream.Send(&pbf.ProcessRequest{Data: &pbf.ProcessRequest_Args{
		Args: args,
	}})
	if err != nil {
		cancel(err)
		return err
	}

	// create_targertMeta() w os.stat (ext + size)
	tMeta, err := createTargertMeta(blender)
	if err != nil {
		cancel(err)
		return err
	}
	err = stream.Send(&pbf.ProcessRequest{Data: &pbf.ProcessRequest_Meta{
		Meta: tMeta,
	}})
	if err != nil {
		cancel(err)
		return err
	}
	// write in blender file (upload targetblob, io.copy to send chunks)
	blenderFileReader, err := os.Open(blender)
	if err != nil {
		cancel(err)
		return err
	}

	errorC := make(chan error, 1+len(args.ExtensionList))

	log.Print("Starting go routine for Upload")
	go loopUploadBlob(errorC, blenderFileReader, writeStream{i: 0, stream: stream})

	dataM := make(map[string]io.WriteCloser)

	for _, ext := range extenstionList {
		dataM[ext], err = os.Create(outfilePath(outDir, ext))
		if err != nil {
			return err
		}
	}
	logC := make(chan Log, 1_000)
	go loopLog(ctx, logC)

	return handleStream(ctx, cancel, stream, dataM, logC)
}

func loopLog(ctx context.Context, logC <-chan Log) {
	doneC := ctx.Done()
out:
	for {
		select {
		case <-doneC:
			break out
		case x := <-logC:
			log.Printf("server log: %s", x)
		}
	}
}

type writeStream struct {
	i      int
	stream pbf.JobManager_ProcessClient
}

func (ws writeStream) Write(p []byte) (n int, err error) {
	log.Printf("write - %d", len(p))

	blob := new(pbf.TargetBlob)
	blob.Data = make([]byte, len(p))
	n = copy(blob.Data, p)
	err = ws.stream.Send(&pbf.ProcessRequest{
		Data: &pbf.ProcessRequest_Blob{
			Blob: blob,
		},
	})
	return
}

func createArgs(extensionList []string) *pbf.ProcessArgs {
	list := make([]string, len(extensionList))
	copy(list, extensionList)
	// Create a new ProcessArgs message using the provided extension list
	args := &pbf.ProcessArgs{
		ExtensionList: list,
	}

	return args
}

func createTargertMeta(blender string) (*pbf.TargetMeta, error) {
	// Create a new ProcessArgs message using the provided extension list
	info, err := os.Stat(blender)
	if err != nil {
		return nil, err
	}
	extension := strings.TrimPrefix(filepath.Ext(blender), ".")
	tMeta := &pbf.TargetMeta{
		Size:      uint64(info.Size()),
		Extension: extension,
	}

	return tMeta, nil
}

func loopUploadBlob(
	errorC chan<- error,
	reader io.ReadCloser,
	writer io.Writer,
) {
	defer reader.Close()
	_, err := io.Copy(writer, reader)
	// we must inform the Process function when we finish our task
	if err != nil {
		errorC <- err
	}
	log.Printf("client upload: %s", err)
}

type Log struct {
	Out string
}

func handleStream(
	ctx context.Context,
	cancel context.CancelCauseFunc,
	stream pbf.JobManager_ProcessClient,
	dataM map[string]io.WriteCloser,
	logC chan<- Log,
) error {
	doneC := ctx.Done()
	defer stream.CloseSend()
	var err error
out:
	for 0 < len(dataM) {
		var msg *pbf.ProcessResponse
		log.Print("client stream - 1")
		msg, err = stream.Recv()
		log.Print("client stream - 2")
		if err == io.EOF {
			err = nil
			break out
		} else if err != nil {
			break out
		}
		switch msg.Data.(type) {
		case *pbf.ProcessResponse_Blob:
			err = handleMessageBlob(msg.GetBlob(), dataM)
		case *pbf.ProcessResponse_Log:
			l := msg.GetLog()
			select {
			case <-doneC:
				break out
			case logC <- Log{
				Out: l.Log,
			}:
			}
		default:
			err = errors.New("unknown message")
		}
		if err != nil {
			break out
		}

	}
	for ext, writer := range dataM {
		// do EOF to file handles
		log.Printf("sending EOF to %s", ext)
		writer.Close()
	}
	log.Printf("1 - finished manager: %s and %d", err, len(dataM))
	return err
}

func handleMessageBlob(
	msg *pbf.TargetBlob,
	dataM map[string]io.WriteCloser,
) error {
	var err error
	if msg.Data == nil {
		msg.Data = []byte{}
	}
	writer, present := dataM[msg.Extension]
	if !present && 0 < len(msg.Data) {
		return fmt.Errorf("unknown output extension %s; %d", msg.Extension, len(msg.Data))

	} else if !present {
		return nil
	}
	if 0 < len(msg.Data) {
		log.Print("client stream - 3")
		_, err = io.Copy(writer, bytes.NewBuffer(msg.Data))
		log.Print("client stream - 4")
		log.Printf("client writing %s %d", msg.Extension, len(msg.Data))
		if err != nil {
			return err
		}
	} else {
		delete(dataM, msg.Extension)
		writer.Close()
	}
	return nil
}
