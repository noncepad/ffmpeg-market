package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	pbf "gitlab.noncepad.com/naomiyoko/ffmpeg-market/proto/ffmpeg"
	"google.golang.org/grpc"
)

type external struct {
	pbf.UnimplementedJobManagerServer
	manager Manager
}

func (e external) Process(stream pbf.JobManager_ProcessServer) error {
	ctx3 := stream.Context()
	go func() {
		<-ctx3.Done()
		log.Print("stream context canceled")
	}()
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
	log.Printf("receiving blender file")
	err = readProcess(ctx, &readStream{stream: stream, i: 0, total: tMeta.Size}, sourceFilePath)
	if err != nil {
		return err
	}
	log.Printf("processing blender file")
	// feed sourceFilePath
	// do conversions in the work pool
	// this blocks until we get read file handles from ffmpeg stdout
	readerList, err := e.manager.SendJob(ctx, sourceFilePath, args.ExtensionList)
	if err != nil {
		return err
	}
	// check correspondence
	if len(args.ExtensionList) != len(readerList) {
		for _, v := range readerList {
			os.Remove(v)
		}
		return fmt.Errorf("expected reader list of length %d but got length %d", len(args.ExtensionList), len(readerList))
	}
	outC := make(chan *readSingle, 100)
	for i, outFp := range readerList {
		go loopRead(ctx, outFp, outC, args.ExtensionList[i])
	}

	i := 0
out:
	for i < len(readerList) {
		select {
		case <-doneC:
			break out
		case r := <-outC:
			if r.err != nil {
				break out
			}
			if r.isDone {
				log.Printf("server finished sending ext %s", r.ext)
				err = stream.Send(&pbf.ProcessResponse{Extenstion: r.ext, Data: []byte{}})
				if err != nil {
					break out
				}
				i++
				continue
			}
			err = stream.Send(&pbf.ProcessResponse{Extenstion: r.ext, Data: r.data})
			if err != nil {
				break out
			}
		}
	}

	log.Printf("exiting process: %s", err)
	// the worker stops working at this point as it has finished sending us file handles

	return err
}

type readSingle struct {
	ext    string
	err    error
	data   []byte
	isDone bool
}

func loopRead(
	ctx context.Context,
	outFp string,
	outC chan<- *readSingle,
	ext string,
) {
	doneC := ctx.Done()
	defer os.Remove(outFp)
	reader, err := os.Open(outFp)
	if err != nil {
		select {
		case <-doneC:
		case outC <- &readSingle{
			ext:    ext,
			err:    err,
			isDone: true,
		}:
		}
		return
	}
	defer reader.Close()
	isDone := false
	var n int
	buf := make([]byte, 2048*16)
out:
	for {
		n, err = reader.Read(buf)
		if err == io.EOF {
			err = nil
			break out
		} else if err != nil {
			break out
		}
		d := make([]byte, n)
		copy(d, buf[0:n])
		log.Printf("1 - writing resp %s %d", ext, len(d))
		select {
		case <-doneC:
			isDone = true
			break out
		case outC <- &readSingle{
			err:    err,
			data:   d,
			ext:    ext,
			isDone: false,
		}:
		}
		log.Printf("2 - writing resp %s %d", ext, len(d))
	}
	if !isDone {
		select {
		case <-doneC:
		case outC <- &readSingle{
			ext:    ext,
			err:    err,
			isDone: true,
		}:
		}
	}

	log.Printf("server completed write %s: %s", ext, err)

}

// routine that reads everything (1 stream)
func readProcess(
	ctx context.Context,
	blobReader io.Reader,
	sourceFilePath string,
) error {
	writeFileHandle, err := os.Create(sourceFilePath)
	if err != nil {
		log.Printf("write file handle error: %s", err)
		return fmt.Errorf("error creating named pipe - 2: %s", err)
	}
	defer writeFileHandle.Close()
	_, err = io.Copy(writeFileHandle, blobReader)
	if err != nil {
		log.Printf("copy error: %s", err)
		return fmt.Errorf("error creating named pipe - 3: %s", err)
	}
	return nil
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
	if rs.total <= uint64(rs.i) {
		n = 0
		err = io.EOF
		return
	}

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
	rs.i += copy(p[0:len(blob.Data)], blob.Data)
	if rs.total < uint64(rs.i) {
		err = fmt.Errorf("too many bytes written: %d vs %d", rs.total, rs.i)
	}
	log.Printf("rs.i %d", rs.i)
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
	log.Printf("removing temporary directory: %s", dir)
	os.RemoveAll(dir)
}
