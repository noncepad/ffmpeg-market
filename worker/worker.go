package worker

import (
	"context"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/noncepad/ffmpeg-market/converter"
	"github.com/noncepad/worker-pool/pool"
	log "github.com/sirupsen/logrus"
)

// Request specifies the structure for job requests
type Request struct {
	Blender string
	Out     []string
}

// Job wraps the details necessary to process a single task.
type Job struct {
	Ctx     context.Context // The context for the job, allowing cancellation signal handling.
	Blender string          // Path to the Blender file to be processed.
	Out     []string        // List of desired output file formats.
	ResultC chan<- Result   // Channel to send the result back to the requester.
}

// Holds a list of strings representing files.
type Result struct {
	Reader []string // Contains file paths of converted files (no file handles are included).
}

// Implementation of a worker capable of processing Jobs.
type simpleWorker struct {
	ctx    context.Context     // context for cancellation
	cancel context.CancelFunc  // call when cancelling the worker
	c      converter.Converter // The video conversion interface
	tmp    string              // Directory for temporary files
}

// Creates instance of a new simpleWorker with temporary directory.
// the tmp directory (tmp) must already exist and be exclusive to the worker
func Create(parentCtx context.Context, c converter.Converter, dirTmp string) (pool.Worker[Request, Result], error) {
	os.MkdirAll(dirTmp, 0755)
	e1 := new(simpleWorker)
	e1.ctx, e1.cancel = context.WithCancel(parentCtx)
	e1.c = c
	e1.tmp = dirTmp
	return e1, nil
}

// Terminate the worker by cancelling its context
func (sw *simpleWorker) Close() error {
	sw.cancel()
	return sw.ctx.Err()
}

// When context is cancelled indicates and error with LoopError
func (sw *simpleWorker) CloseSignal() <-chan error {
	signalC := make(chan error, 1)
	go loopError(sw.ctx, signalC)
	return signalC
}

// Listens for the context's done signal and relays the context's error to the provided channel.
func loopError(ctx context.Context, errorC chan<- error) {
	<-ctx.Done()
	errorC <- ctx.Err()
}

// Generates the path to the intermediary, tmp AVI file used during conversion.
func (sw *simpleWorker) intermediaryAviFilePath() string {
	return sw.tmp + "/blah.avi"
}

// Used to pass around results from the conversion go routine.
type convertResult struct {
	i   int
	err error
}

// loopConvert calls convert and sends the result over a channel.
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
	log.Debugf("converting from %s to %s", inFp, outFp)
	cr.err = c.Convert(ctx, inFp, outFp)
	if cr.err != nil {
		ansC <- cr
		return
	}
	log.Debugf("converted from %s to %s", inFp, outFp)
	ansC <- cr
}

// randomString generates a random string of a specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
