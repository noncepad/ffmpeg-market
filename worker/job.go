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
	// call render to convert blender file to .avi format
	err := sw.c.Render(jobCtx, payload.Blender, sw.intermediaryAviFilePath())
	if err != nil {
		return Result{Reader: []string{}}, err // Return empty result and error if render fails.
	}

	// Define an array of file handles (Reader)
	readerList := make([]string, len(payload.Out))
	ansC := make(chan *convertResult, len(readerList)) // Channel to receive results from conversions.
	// call convert to convert avi (in working dir) to mp4, mkv, gif, etc.
	wg := &sync.WaitGroup{}                   // WaitGroup for synchronizing go routines.
	wg.Add(len(payload.Out))                  // Add the number of expected conversions to the WaitGroup.
	ctx, cancel := context.WithCancel(jobCtx) // Create a cancellable sub-context.
	for i, fileExtension := range payload.Out {
		r := randomString(10)                                                                        // Generate a random string for unique file paths.
		readerList[i] = filepath.Join(sw.tmp, "tmpout"+r+"."+strings.TrimPrefix(fileExtension, ".")) // Append generated path to reader list.
		go loopConvert(ctx, wg, sw.c, i, ansC, sw.intermediaryAviFilePath(), readerList[i])          // Start conversion in a separate goroutine.
	}
out:
	for k := 0; k < len(payload.Out); k++ {
		var cr *convertResult
		select {
		case <-doneC:
			err = ctx.Err() // If the job context is done, assign error and break out.
			break out
		case cr = <-ansC: // Receive a conversion result.
			if cr.err != nil {
				err = cr.err // If there is an error in conversion, assign it and break out.
				break out
			}
			// Log successful assignment of the file pointer from the convert result.
			log.Printf("assigning reader cr %d; fp %s", cr.i, readerList[cr.i])
		}
	}
	cancel()  // Cancel the sub-context.
	wg.Wait() // Wait for all goroutines to finish.

	if err != nil {
		for _, outFp := range readerList {
			os.Remove(outFp) // Clean up: remove the output files.
		}
		readerList = nil
		return Result{Reader: []string{}}, nil // Return empty result due to error.
	}
	return Result{Reader: readerList}, nil // Return the result with the list of output files.

}
