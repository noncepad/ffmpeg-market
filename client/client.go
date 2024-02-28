package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	pbf "github.com/noncepad/ffmpeg-market/proto/ffmpeg"
	"google.golang.org/grpc"
)

type Client interface {
	// ProcessRequest is responsible for processing the rendering request.
	ProcessRequest(ctx context.Context, blender string, outDir string, extenstionList []string) error
}

type external struct {
	client pbf.JobManagerClient // JobManagerClient is a gRPC client for the job manager service.
}

// simple implementation of Client
func Create(ctx context.Context, conn *grpc.ClientConn) Client {
	return external{client: pbf.NewJobManagerClient(conn)} // Initializes and returns a new gRPC client
}

func (e external) ProcessRequest(parentCtx context.Context, blender string, outDir string, extenstionList []string) error {
	// Check if the file exists and has the correct extension
	if _, err := os.Stat(blender); os.IsNotExist(err) {
		return fmt.Errorf("processRequest - 1: file does not exist, err: %s", err)
	} else if filepath.Ext(blender) != ".blend" {
		// Ensure the provided path has the correct '.blend' extension for Blender files.
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
		// Validate that each requested video extension is supported.
		if !supportedExtensions[ext] {
			return fmt.Errorf("processRequest - 3: extension %s is not compatible with ffmpeg", ext)
		}
	}
	// set context with cancel
	ctx, cancel := context.WithCancelCause(parentCtx)

	// Start a streaming gRPC call to process the input file.
	stream, err := e.client.Process(ctx)
	if err != nil {
		cancel(err)
		return err
	}
	// Prepare arguments for the gRPC call based on provided file extensions.
	args := createArgs(extenstionList)
	err = stream.Send(&pbf.ProcessRequest{Data: &pbf.ProcessRequest_Args{
		Args: args,
	}})
	if err != nil {
		cancel(err)
		return err
	}
	// Collect metadata about the source Blender file (ext + size)
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
	// Open the Blender file for reading.
	blenderFileReader, err := os.Open(blender)
	if err != nil {
		cancel(err)
		return err
	}

	errorC := make(chan error, 1+len(args.ExtensionList)) // Channel to communicate errors from concurrent operations.

	log.Print("Starting go routine for Upload")
	go loopUploadBlob(errorC, blenderFileReader, writeStream{i: 0, stream: stream}) // Start uploading the Blender file in a separate goroutine.

	dataM := make(map[string]io.WriteCloser) // Map to maintain handles to output files for each extension.

	for _, ext := range extenstionList {
		// Create an output file for each requested extension.
		dataM[ext], err = os.Create(outfilePath(outDir, ext))
		if err != nil {
			return err
		}
	}
	logC := make(chan Log, 1_000) // Log channel to capture and handle logs.
	go loopLog(ctx, logC)         // Start handling logs in a separate goroutine.

	// Handle incoming messages from the server and manage output file writing.
	return handleStream(ctx, cancel, stream, dataM, logC)
}

func createArgs(extensionList []string) *pbf.ProcessArgs {
	list := make([]string, len(extensionList))
	copy(list, extensionList) // Duplicate the extension list to avoid modifying the original slice.
	// Create a new ProcessArgs message using the provided extension list
	args := &pbf.ProcessArgs{
		ExtensionList: list,
	}

	return args
}

func createTargertMeta(blender string) (*pbf.TargetMeta, error) {
	info, err := os.Stat(blender) // Get file information for the Blender file.
	if err != nil {
		return nil, err
	}
	extension := strings.TrimPrefix(filepath.Ext(blender), ".") // Trim the leading dot from the file extension.
	tMeta := &pbf.TargetMeta{
		Size:      uint64(info.Size()), // Set the size of the Blender file.
		Extension: extension,           // Set the file extension.
	}

	return tMeta, nil
}

// upload the target file to the server
func loopUploadBlob(
	errorC chan<- error,
	reader io.ReadCloser,
	writer io.Writer,
) {
	defer reader.Close()              // Ensure the file reader is closed at the end of this function.
	_, err := io.Copy(writer, reader) // Stream the entire Blender file to the server.
	// we must inform the Process function when we finish our task
	if err != nil {
		errorC <- err
	}
	log.Printf("client upload: %s", err) // Log the result of the upload operation.
}

func outfilePath(dir string, ext string) string {
	// Construct an output file path based on the output directory and file extension
	return fmt.Sprintf("%s/out.%s", dir, ext)
}
