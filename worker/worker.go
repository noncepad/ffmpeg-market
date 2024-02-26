package worker

import (
	"context"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/noncepad/worker-pool/pool"
	log "github.com/sirupsen/logrus"
	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/converter"
)

type Request struct {
	Blender string
	Out     []string
}

type Job struct {
	Ctx     context.Context
	Blender string   // filepath
	Out     []string // list of file extensions
	ResultC chan<- Result
}

type Result struct {
	Reader []string // no file handle
}

type simpleWorker struct {
	ctx    context.Context
	cancel context.CancelFunc
	c      converter.Converter
	tmp    string //directory
}

// the tmp directory (tmp) must already exist and be exclusive to the worker
func Create(parentCtx context.Context, c converter.Converter, dirTmp string) (pool.Worker[Request, Result], error) {
	os.MkdirAll(dirTmp, 0755)
	e1 := new(simpleWorker)
	e1.ctx, e1.cancel = context.WithCancel(parentCtx)
	e1.c = c
	e1.tmp = dirTmp
	return e1, nil
}

func (sw *simpleWorker) Close() error {
	sw.cancel()
	return sw.ctx.Err()
}
func (sw *simpleWorker) CloseSignal() <-chan error {
	signalC := make(chan error, 1)
	go loopError(sw.ctx, signalC)
	return signalC
}

func loopError(ctx context.Context, errorC chan<- error) {
	<-ctx.Done()
	errorC <- ctx.Err()
}

func (sw *simpleWorker) intermediaryAviFilePath() string {
	return sw.tmp + "/blah.avi"
}

type convertResult struct {
	i   int
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
	log.Debugf("converting from %s to %s", inFp, outFp)
	cr.err = c.Convert(ctx, inFp, outFp)
	if cr.err != nil {
		ansC <- cr
		return
	}
	log.Debugf("converted from %s to %s", inFp, outFp)
	ansC <- cr
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
