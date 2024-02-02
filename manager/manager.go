package manager

import (
	"context"
	"fmt"

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
}

func Create(parentCtx context.Context, conf *Configuration, workerCb WorkerCallback) (Manager, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	// no buffer so program will block if fails
	// how many jobs do you want in the buffer? for now no buffer
	jobC := make(chan worker.Job)
	// creating worker pool
	for i := 0; i < n; i++ {
		w, err := workerCb(ctx, conf)
		if err != nil {
			cancel()
			return Manager{}, fmt.Errorf("Failed to create worker: ", err)
		} else {
			fmt.Printf("Worker created successfully")
		}
		go loopDoJob(w, jobC, cancel)
	}
	return Manager{}, nil

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
	fmt.Printf("Closing loop:", err)
}
