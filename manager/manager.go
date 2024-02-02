package manager

import (
	"context"
	"fmt"
	"io"

	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/worker"
)

type Configuration struct {
	DirWork    string
	BinFfmpeg  string
	BinBlender string
	ListenUrl  string
	MaxJobs    int
}

type WorkerCallback func(ctx context.Context, conf *Configuration) (worker.Worker, error)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	jobC   chan worker.Job
}

func Create(parentCtx context.Context, conf *Configuration, workerCb WorkerCallback) (Manager, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	// no buffer so program will block if fails
	// how many jobs do you want in the buffer? for now no buffer
	jobC := make(chan worker.Job)
	// creating worker pool
	n := conf.MaxJobs
	for i := 0; i < n; i++ {
		w, err := workerCb(ctx, conf)
		if err != nil {
			cancel()
			return Manager{}, fmt.Errorf("failed to create worker: %s", err)
		} else {
			fmt.Printf("Worker created successfully")
		}
		go loopDoJob(w, jobC, cancel)
	}
	return Manager{
		ctx:    ctx,
		cancel: cancel,
		jobC:   jobC,
	}, nil

}

// jobC is a read only channel
// whenever you have a cancel, do defer cancel
// this function will block forever until the program dies
func loopDoJob(w worker.Worker, jobC <-chan worker.Job, cancel context.CancelFunc) {
	defer cancel()
	var job worker.Job
	var err error
	signalC := w.CloseSignal()
out:
	for {
		select {
		// we read from signalC and get an error and put t in the err variable
		case err = <-signalC:

			break out
		// we read from jobC and get a job and put it in the job variable
		case job = <-jobC:
			result := w.Run(job)
			// ResultC always has a buffer of 1, we will not block
			job.ResultC <- result
		}
	}
	fmt.Printf("closing loop:", err)
}

func (m Manager) SendJob(
	ctx context.Context,
	Blender string,
	Out []string,
) ([]io.ReadCloser, error) {
	ResultC := make(chan worker.Result)
	job := worker.Job{
		Ctx:     ctx,     //context.Context
		Blender: Blender, //string   // filepath
		Out:     Out,     //[]string // list of file extensions
		ResultC: ResultC, //chan<- Result
	}
	select {
	case m.jobC <- job:
		// Job was successfully sent to channel
	case <-m.ctx.Done():
		// The manager's context has been cancelled, return an error
		return nil, m.ctx.Err()
	case <-job.Ctx.Done():
		// The job's context was been cancelled, return an error
		return nil, job.Ctx.Err()
	}
	select {
	case result := <-ResultC:
		// Job was successfully sent to channel
		return result.Reader, nil
	case <-m.ctx.Done():
		// The manager's context has been cancelled, return an error
		return nil, m.ctx.Err()
	case <-job.Ctx.Done():
		// The job's context was been cancelled, return an error
		return nil, job.Ctx.Err()

	}
}
