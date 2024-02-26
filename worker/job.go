package worker

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/noncepad/worker-pool/pool"
)

func (sw *simpleWorker) Run(job pool.Job[Request]) (Result, error) {
	jobCtx := job.Ctx()
	doneC := jobCtx.Done()
	payload := job.Payload()
	// call render to convert blender file to avi
	err := sw.c.Render(jobCtx, payload.Blender, sw.intermediaryAviFilePath())
	if err != nil {
		return Result{Reader: []string{}}, err
	}

	// Define an array of file handles (Reader)
	readerList := make([]string, len(payload.Out))
	ansC := make(chan *convertResult, len(readerList))
	// call convert to convert avi (in working dir) to mp4, mkv, gif
	wg := &sync.WaitGroup{}
	wg.Add(len(payload.Out))
	ctx, cancel := context.WithCancel(jobCtx)
	for i, fileExtension := range payload.Out {
		r := randomString(10)
		readerList[i] = filepath.Join(sw.tmp, "tmpout"+r+"."+strings.TrimPrefix(fileExtension, "."))
		go loopConvert(ctx, wg, sw.c, i, ansC, sw.intermediaryAviFilePath(), readerList[i])
	}
out:
	for k := 0; k < len(payload.Out); k++ {
		var cr *convertResult
		select {
		case <-doneC:
			err = ctx.Err()
			break out
		case cr = <-ansC:
			if cr.err != nil {
				err = cr.err
				break out
			}
			log.Printf("assigning reader cr %d; fp %s", cr.i, readerList[cr.i])
		}
	}
	cancel()
	wg.Wait()
	if err != nil {
		for _, outFp := range readerList {
			os.Remove(outFp)
		}
		readerList = nil
		return Result{Reader: []string{}}, nil
	}
	return Result{Reader: readerList}, nil
}
