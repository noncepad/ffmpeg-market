package worker

import (
	"context"
	"io"
	"os"

	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/converter"
)

type Job struct {
	Ctx     context.Context
	Blender string   // filepath
	Out     []string // list of file extensions
	ResultC chan<- Result
}

type Result struct {
	Reader []io.ReadCloser // file handle
	Err    error
}

// err: if render fails -> no file handle for conversion
type Worker interface {
	Run(Job) Result
	// this is ther workers version of ctx.Done()
	CloseSignal() <-chan error
}

type simpleWorker struct {
	ctx context.Context
	c   converter.Converter
	tmp string //directory
}

// the tmp directory (tmp) must already exist and be exclusive to the worker
func Create(ctx context.Context, c converter.Converter, dirTmp string) (Worker, error) {
	os.MkdirAll(dirTmp, 0755)
	e1 := new(simpleWorker)
	e1.ctx = ctx
	e1.c = c
	e1.tmp = dirTmp
	return e1, nil
}

func (sw *simpleWorker) intermediaryAviFilePath() string {
	return sw.tmp + "/blah.avi"
}

func (sw *simpleWorker) Run(job Job) Result {
	// call render to convert blender file to avi
	err := sw.c.Render(job.Ctx, job.Blender, sw.intermediaryAviFilePath())
	if err != nil {
		return Result{Reader: nil, Err: err}
	}

	// Define an array of file handles (Reader)
	readerList := make([]io.ReadCloser, len(job.Out))

	// call convert to convert avi (in working dir) to mp4, mkv, gif
	for i, fileExtension := range job.Out {
		fileHandle, err := sw.c.Convert(job.Ctx, sw.intermediaryAviFilePath(), fileExtension)
		if err != nil {
			// if an error occurs, close all handles and retun err
			for _, reader := range readerList {
				if reader != nil {
					reader.Close()
				}
			}
			return Result{Reader: nil, Err: err}
		}
		readerList[i] = fileHandle
	}
	return Result{Reader: readerList, Err: nil}

}

func (sw *simpleWorker) CloseSignal() <-chan error {
	// Create a channel to receive signals
	signalC := make(chan error)

	// Run a goroutine to monitor for context cancellation or errors
	go func() {
		// Block until the worker's context is cancelled or an error occurs
		<-sw.ctx.Done()

		// Send a signal to the channel if to close
		signalC <- sw.ctx.Err()
	}()

	// Return the channel for receiving signals
	return signalC
}
