package worker

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	doneC := job.Ctx.Done()
	// call render to convert blender file to avi
	err := sw.c.Render(job.Ctx, job.Blender, sw.intermediaryAviFilePath())
	if err != nil {
		return Result{Reader: nil, Err: err}
	}

	// Define an array of file handles (Reader)
	readerList := make([]io.ReadCloser, len(job.Out))
	ansC := make(chan *convertResult, len(readerList))
	// call convert to convert avi (in working dir) to mp4, mkv, gif
	wg := &sync.WaitGroup{}
	wg.Add(len(job.Out))
	ctx, cancel := context.WithCancel(job.Ctx)
	for i, fileExtension := range job.Out {
		outFp := filepath.Join(sw.tmp, "tmpout."+strings.TrimPrefix(fileExtension, "."))
		go loopConvert(ctx, wg, sw.c, i, ansC, sw.intermediaryAviFilePath(), outFp)
	}
out:
	for k := 0; k < len(job.Out); k++ {
		var cr *convertResult
		select {
		case <-doneC:
			err = job.Ctx.Err()
			break out
		case cr = <-ansC:
			if cr.err != nil {
				break out
			}
			outFp := filepath.Join(sw.tmp, "tmpout."+strings.TrimPrefix(job.Out[cr.i], "."))
			log.Printf("assigning reader cr %d; fp %s; reader %s", cr.i, outFp, cr.r)
			readerList[cr.i] = cr.r
		}
	}
	cancel()
	wg.Wait()
	if err != nil {
		for _, reader := range readerList {
			if reader != nil {
				reader.Close()
			}
		}
		readerList = nil
	}
	return Result{Reader: readerList, Err: err}
}

type convertResult struct {
	i   int
	r   io.ReadCloser
	err error
}

func loopConvert(
	ctx context.Context,
	wg *sync.WaitGroup,
	c converter.Converter,
	i int,
	ansC chan<- *convertResult,
	inFp string,
	outFp string,

) {
	defer wg.Done()
	cr := new(convertResult)
	cr.i = i
	cr.err = c.Convert(ctx, inFp, outFp)
	if cr.err != nil {
		ansC <- cr
		return
	}
	cr.r, cr.err = os.Open(outFp)
	ansC <- cr
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
