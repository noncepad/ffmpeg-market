package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	pbf "gitlab.noncepad.com/naomiyoko/ffmpeg-market/proto/ffmpeg"
	"google.golang.org/grpc"
)

type Client interface {
	ProcessRequest(ctx context.Context, blender string, outDir string, extenstionList []string) error
}

type external struct {
	client pbf.JobManagerClient
}

// simple implementation of Client
func Create(ctx context.Context, conn *grpc.ClientConn) Client {
	// Create a new gRPC client
	return external{client: pbf.NewJobManagerClient(conn)}
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

// upload the target file to the server
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

func outfilePath(dir string, ext string) string {
	return fmt.Sprintf("%s/out.%s", dir, ext)
}
