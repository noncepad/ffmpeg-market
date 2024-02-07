package client

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	pbf "gitlab.noncepad.com/naomiyoko/ffmpeg-market/proto/ffmpeg"
	"google.golang.org/grpc"
)

type Client interface {
	ProcessRequest(ctx context.Context, blender string, extenstionList []string) ([]io.ReadCloser, error)
}

type external struct {
	client pbf.JobManagerClient
}

func Create(ctx context.Context, conn *grpc.ClientConn) Client {
	// Create a new gRPC client
	return external{client: pbf.NewJobManagerClient(conn)}
}

func (e external) ProcessRequest(parentCtx context.Context, blender string, extenstionList []string) ([]io.ReadCloser, error) {
	// set context with cancel
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	// gRPC call to the server
	stream, err := e.client.Process(ctx)
	if err != nil {
		return nil, err
	}
	// create_args() w extensionlist
	args := createArgs(extenstionList)
	err = stream.Send(&pbf.ProcessRequest{Data: &pbf.ProcessRequest_Args{
		Args: args,
	}})
	if err != nil {
		return nil, err
	}

	// create_targertMeta() w os.stat (ext + size)
	tMeta, err := createTargertMeta(blender)
	if err != nil {
		return nil, err
	}
	err = stream.Send(&pbf.ProcessRequest{Data: &pbf.ProcessRequest_Meta{
		Meta: tMeta,
	}})
	if err != nil {
		return nil, err
	}
	// write in blender file (upload targetblob, io.copy to send chunks)
	blenderFileReader, err := os.Open(blender)
	if err != nil {
		return nil, err
	}

	errorC := make(chan error, 1+len(args.ExtensionList))
	go loopUploadBlob(errorC, blenderFileReader, writeStream{i: 0, stream: stream})

	// create array to return readers
	readerList := make([]io.ReadCloser, len(extenstionList))
	dataM := make(map[string]chan<- []byte)
	for i, ext := range extenstionList {
		dataC := make(chan []byte, 1000)
		dataM[ext] = dataC
		readerList[i] = &readStream{ctx: ctx, i: 0, dataC: dataC, cancel: cancel, ext: ext}
	}

	go loopManager(ctx, cancel, errorC, stream, dataM)

	// return reader array
	return readerList, nil
}

type writeStream struct {
	i      int
	stream pbf.JobManager_ProcessClient
}

func (ws writeStream) Write(p []byte) (n int, err error) {
	data := make([]byte, len(p))

	resp := new(pbf.ProcessResponse)
	resp.Data = data
	n = copy(data, p)
	ws.i += n
	//n = len(resp.Data)
	return
}

type readStream struct {
	ctx    context.Context
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
		}
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
	extension := filepath.Ext(blender)
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
	errorC <- err
}

func loopManager(ctx context.Context, cancel context.CancelFunc, errorC chan<- error, stream pbf.JobManager_ProcessClient, dataM map[string]chan<- []byte) {
	defer cancel()
	doneC := ctx.Done()
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			err = nil
			for _, dataC := range dataM {
				select {
				case <-doneC:
					return
				case dataC <- []byte{}:
				}
			}
			return
		} else if err != nil {
			return
		}
		dataC, present := dataM[msg.Extenstion]
		if !present {
			return
		}
		select {
		case <-doneC:
			return
		case dataC <- msg.Data:
		}

	}

}

func (rs readStream) Close() error {
	rs.cancel()
	return nil
}
