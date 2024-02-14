package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	pbf "gitlab.noncepad.com/naomiyoko/ffmpeg-market/proto/ffmpeg"
)

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
