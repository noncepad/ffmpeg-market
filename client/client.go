package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	pbf "gitlab.noncepad.com/naomiyoko/ffmpeg-market/proto/ffmpeg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	ProcessRequest(ctx context.Context, blender string, extenstionList []string) ([]io.ReadCloser, error)
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
	doneC := ctx.Done()
	conn, err := grpc.DialContext(ctx, args[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client := Create(ctx, conn)
	blenderIn := args[1]
	dir := args[2]
	extList := args[3:]
	readerList, err := client.ProcessRequest(ctx, blenderIn, extList)
	if err != nil {
		return err
	}
	errorC := make(chan error, len(readerList))
	for i, reader := range readerList {
		f, err := os.Create(fmt.Sprintf("%s/out.%s", dir, extList[i]))
		if err != nil {
			return err
		}
		log.Printf("Run - 2: Starting goroutine for writing %s\n", f.Name())
		go loopWriteFile(ctx, errorC, f, reader)
	}
out:
	for i := 0; i < len(readerList); i++ {
		select {
		case <-doneC:
			err = ctx.Err()
			break out
		case err = <-errorC:
			if err != nil {
				break out
			}
		}
	}
	return err
}

func loopWriteFile(ctx context.Context, errorC chan<- error, fileOut io.WriteCloser, reader io.ReadCloser) {
	_, err := io.Copy(fileOut, reader)
	log.Printf("finished write-file copy: %s", err)
	select {
	case <-ctx.Done():
	case errorC <- err:

	}

}

func Create(ctx context.Context, conn *grpc.ClientConn) Client {
	// Create a new gRPC client
	return external{client: pbf.NewJobManagerClient(conn)}
}

func (e external) ProcessRequest(parentCtx context.Context, blender string, extenstionList []string) ([]io.ReadCloser, error) {
	// Check if the file exists and has the correct extension
	if _, err := os.Stat(blender); os.IsNotExist(err) {
		return nil, fmt.Errorf("processRequest - 1: file does not exist, err: %s", err)
	} else if filepath.Ext(blender) != ".blend" {
		return nil, fmt.Errorf("process Request - 2: file must have .blend extension, err %s", err)
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
			return nil, fmt.Errorf("processRequest - 3: extension %s is not compatible with ffmpeg", ext)
		}
	}
	// set context with cancel
	ctx, cancel := context.WithCancel(parentCtx)
	wg := &sync.WaitGroup{}

	// gRPC call to the server
	stream, err := e.client.Process(ctx)
	if err != nil {
		cancel()
		return nil, err
	}
	// create_args() w extensionlist
	args := createArgs(extenstionList)
	err = stream.Send(&pbf.ProcessRequest{Data: &pbf.ProcessRequest_Args{
		Args: args,
	}})
	if err != nil {
		cancel()
		return nil, err
	}

	// create_targertMeta() w os.stat (ext + size)
	tMeta, err := createTargertMeta(blender)
	if err != nil {
		cancel()
		return nil, err
	}
	err = stream.Send(&pbf.ProcessRequest{Data: &pbf.ProcessRequest_Meta{
		Meta: tMeta,
	}})
	if err != nil {
		cancel()
		return nil, err
	}
	// write in blender file (upload targetblob, io.copy to send chunks)
	blenderFileReader, err := os.Open(blender)
	if err != nil {
		cancel()
		return nil, err
	}

	errorC := make(chan error, 1+len(args.ExtensionList))

	log.Print("Starting go routine for Upload")
	go loopUploadBlob(errorC, blenderFileReader, writeStream{i: 0, stream: stream})

	// create array to return readers
	readerList := make([]io.ReadCloser, len(extenstionList))
	dataM := make(map[string]chan<- []byte)
	wg.Add(len(extenstionList))

	for i, ext := range extenstionList {
		dataC := make(chan []byte, 1000)
		dataM[ext] = dataC
		readerList[i] = &readStream{ctx: ctx, wg: wg, i: 0, dataC: dataC, cancel: cancel, ext: ext}
	}

	go loopManager(ctx, errorC, stream, dataM)
	go loopWait(wg, cancel)
	go loopStop(ctx, cancel, errorC)
	// return reader array
	return readerList, nil
}
func loopStop(ctx context.Context, cancel context.CancelFunc, errorC <-chan error) {
	defer cancel()
	var err error
	select {
	case <-ctx.Done():
	case err = <-errorC:
	}
	log.Printf("stop error: %s", err)
}
func loopWait(wg *sync.WaitGroup, cancel context.CancelFunc) {
	wg.Wait()
	cancel()
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

type readStream struct {
	ctx    context.Context
	wg     *sync.WaitGroup
	i      int
	dataC  <-chan []byte
	cancel context.CancelFunc
	ext    string
}

// we need to put data into the p []byte array.
func (rs readStream) Read(p []byte) (n int, err error) {
	select {
	case <-rs.ctx.Done():
		err = rs.ctx.Err()
	case d := <-rs.dataC:
		if len(d) > len(p) {
			err = fmt.Errorf("buffer to small: %d vs. %d", len(d), len(p))
		} else {
			n = copy(p, d)
			if n == 0 {
				err = io.EOF
			}
		}
	}
	if err == io.EOF {
		rs.wg.Done()
	} else if err != nil {
		rs.cancel()
		rs.wg.Done()
	}
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

func loopManager(ctx context.Context, errorC chan<- error, stream pbf.JobManager_ProcessClient, dataM map[string]chan<- []byte) {

	doneC := ctx.Done()
	var err error
out:
	for 0 < len(dataM) {
		var msg *pbf.ProcessResponse
		msg, err = stream.Recv()
		if err == io.EOF {
			err = nil
			break out
		} else if err != nil {
			break out
		}
		dataC, present := dataM[msg.Extenstion]
		if !present {
			err = fmt.Errorf("unknown output extension %s", msg.Extenstion)
			break out
		}
		select {
		case <-doneC:
			break out
		case dataC <- msg.Data:
			log.Printf("mgr send (ext %s) %d", msg.Extenstion, len(msg.Data))
			if len(msg.Data) == 0 {
				delete(dataM, msg.Extenstion)
			}
		}
	}

	if err != nil {
		errorC <- err
	}

	log.Printf("finished manager: %s and %d", err, len(dataM))
}

func (rs readStream) Close() error {
	rs.cancel()
	return nil
}
