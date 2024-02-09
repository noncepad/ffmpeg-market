package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	pbf "gitlab.noncepad.com/naomiyoko/ffmpeg-market/proto/ffmpeg"
	"google.golang.org/grpc"
)

type external struct {
	pbf.UnimplementedJobManagerServer
	manager Manager
}

func (e external) Process(stream pbf.JobManager_ProcessServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	doneC := ctx.Done()
	defer cancel()
	/*
	   message ProcessRequest{
	     oneof data{
	       ProcessArgs args = 1;
	       TargetMeta meta = 2;
	       TargetBlob blob = 3;
	     }
	   }
	*/
	// assume the user first sends us ProcessArgs; fail otherwise
	args, err := process_readArgs(stream)
	if err != nil {
		return err
	}

	// second, read TargetMeta; fail otherwise
	tMeta, err := process_readTargetMeta(stream)
	if err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("/tmp", "blender*")
	if err != nil {
		return err
	}
	go loopCleanUpDir(ctx, tmpDir)

	// TODO: clean up tmpdir
	sourceFilePath := tmpDir + "/blender." + tMeta.Extension // we need to think of a temporary directory here

	// we need to read in the blender file (only TargetBlob)
	// 1 error for go routine we spawn here
	errorC := make(chan error, 1+len(args.ExtensionList))
	go loopReadProcess(ctx, &readStream{stream: stream, i: 0, total: tMeta.Size}, errorC, sourceFilePath)

	// feed sourceFilePath
	// do conversions in the work pool
	// this blocks until we get read file handles from ffmpeg stdout
	readerList, err := e.manager.SendJob(ctx, sourceFilePath, args.ExtensionList)
	if err != nil {
		return err
	}
	// check correspondence
	if len(args.ExtensionList) != len(readerList) {
		return fmt.Errorf("expected reader list of length %d but got length %d", len(args.ExtensionList), len(readerList))
	}

	/*
			message ProcessResponse{
		    string extenstion=1;
		    bytes data=2;
		}
	*/
	// there is a one to one correspondence between args.ExtensionList and readerList
	for i, ext := range args.ExtensionList {
		// read stdout from the ffmpeg process in manager, and send the data to the customer
		go loopProcessWrite(errorC, readerList[i], &writeStream{i: 0, ext: ext, stream: stream})
	}

	// wait until the 1+n go routines above finish
	for i := 0; i < 1+len(args.ExtensionList); i++ {
		select {
		case <-doneC:
			return ctx.Err()
		case err = <-errorC:
			if err != nil {
				return err
			}
		}
	}
	// the worker stops working at this point as it has finished sending us file handles

	return nil
}

// go routine that reads everything (1 stream)
func loopReadProcess(
	ctx context.Context,
	blobReader io.Reader,
	errorC chan<- error,
	sourceFilePath string,
) {
	doneC := ctx.Done()
	// Create the named pipe

	writeFileHandle, err := os.Create(sourceFilePath)
	if err != nil {
		select {
		case <-doneC:
		case errorC <- fmt.Errorf("error creating named pipe - 2: %s", err):
		}
		return
	}
	defer writeFileHandle.Close()
	_, err = io.Copy(writeFileHandle, blobReader)
	if err != nil {
		select {
		case <-doneC:
		case errorC <- fmt.Errorf("error creating named pipe - 3: %s", err):
		}
	} else {
		// tell the Process function we have completed our task
		select {
		case <-doneC:
		case errorC <- nil:
		}
	}
}

func loopProcessWrite(
	errorC chan<- error,
	reader io.ReadCloser,
	writer io.Writer,
) {
	defer reader.Close()
	_, err := io.Copy(writer, reader)
	// we must inform the Process function when we finish our task
	errorC <- err
}

// implements io.Writer interface
type writeStream struct {
	ext    string
	i      int
	stream pbf.JobManager_ProcessServer
}

// the data in already in p []byte; we need to Send via stream
// n is the number of bytes in the array
func (ws *writeStream) Write(p []byte) (n int, err error) {
	data := make([]byte, len(p))

	resp := new(pbf.ProcessResponse)
	resp.Extenstion = ws.ext
	resp.Data = data
	n = copy(data, p)
	ws.i += n
	//n = len(resp.Data)
	return
}

// TODO: sanitize args
// grab the processArg data from user
func process_readArgs(
	stream pbf.JobManager_ProcessServer,
) (*pbf.ProcessArgs, error) {
	msg, err := stream.Recv()
	if err != nil {
		return nil, err
	}
	// from the protobuf, we have: "oneof data"
	switch msg.Data.(type) {
	case *pbf.ProcessRequest_Args:
		return msg.GetArgs(), nil
	case *pbf.ProcessRequest_Meta:
		return nil, errors.New("received meta, not args")
	case *pbf.ProcessRequest_Blob:
		return nil, errors.New("received blob, not args")
	default:
		return nil, errors.New("received unknown, not args")
	}
}

// TODO: sanitize tMeta
func process_readTargetMeta(stream pbf.JobManager_ProcessServer) (*pbf.TargetMeta, error) {
	msg, err := stream.Recv()
	if err != nil {
		return nil, err
	}
	// from the protobuf, we have: "oneof data"
	switch msg.Data.(type) {
	case *pbf.ProcessRequest_Args:
		return nil, errors.New("received args, not meta")
	case *pbf.ProcessRequest_Meta:
		return msg.GetMeta(), nil
	case *pbf.ProcessRequest_Blob:
		return nil, errors.New("received blob, not meta")
	default:
		return nil, errors.New("received unknown, not meta")
	}

}

// implements io.Reader interface
type readStream struct {
	i      int
	total  uint64
	stream pbf.JobManager_ProcessServer
}

// we need to put data into the p []byte array.
func (rs *readStream) Read(p []byte) (n int, err error) {
	n = 0
	msg, err := rs.stream.Recv()
	if err != nil {
		return
	}
	var blob *pbf.TargetBlob
	// from the protobuf, we have: "oneof data"
	switch msg.Data.(type) {
	case *pbf.ProcessRequest_Args:
		err = errors.New("received args, not blob")
	case *pbf.ProcessRequest_Meta:
		err = errors.New("received blob, not blob")
	case *pbf.ProcessRequest_Blob:
		blob = msg.GetBlob()
	default:
		err = errors.New("received unknown, not blob")
	}
	if err != nil {
		return
	}

	// the data we need is in the "blob" message;
	// we need to get that data into "p"
	if blob.Data == nil {
		err = errors.New("data is blank")
		return
	}
	n = len(blob.Data)
	if len(p) < n {
		err = fmt.Errorf("buffer is too small: %d vs %d", len(p), n)
		return
	}
	// write down how many bytes we have written so far
	rs.i += copy(p, blob.Data)
	if rs.total < uint64(rs.i) {
		err = fmt.Errorf("too many bytes written: %d vs %d", rs.total, rs.i)
	}
	return
}

func (m Manager) Add(s *grpc.Server) {
	e1 := external{
		manager: m,
	}
	pbf.RegisterJobManagerServer(s, e1)
}

// when the request is finished, clean up with this function
func loopCleanUpDir(
	ctx context.Context,
	dir string,
) {
	// delete the temporary directory
	<-ctx.Done()
	os.RemoveAll(dir)
}
