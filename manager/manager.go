package manager

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/converter"
	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/worker"
	"google.golang.org/grpc"
)

type Configuration struct {
	DirWork    string
	BinFfmpeg  string
	BinBlender string
	ListenUrl  string
	MaxJobs    int
}

func (conf *Configuration) Check() error {
	if conf.DirWork == "" {
		return fmt.Errorf("DirWork not set")
	}
	_, err := os.Stat(conf.DirWork)
	if err != nil {
		return fmt.Errorf("DirWork directory does not exist: %v", err)
	}
	if conf.BinFfmpeg == "" {
		return fmt.Errorf("BinFfmpeg not set")
	}
	if err := ExecCheck(conf.BinBlender); err != nil {
		fmt.Println(err)
	}
	if err := ExecCheck(conf.BinFfmpeg); err != nil {
		fmt.Println(err)
	}
	if conf.ListenUrl == "" {
		return fmt.Errorf("listenUrl not set")
	}
	// expected number of jobs is from range 1 till number of workers
	if conf.MaxJobs <= 0 {
		return fmt.Errorf("MaxJobs has unexpected value")
	}
	return nil
}
func DefaultCreateWorker(ctx context.Context, conf *Configuration, i int) (worker.Worker, error) {
	converterConfig := &converter.Configuration{
		BinFfmpeg:  conf.BinFfmpeg,
		BinBlender: conf.BinBlender,
	}

	conv, err := converter.CreateSimpleConverter(ctx, converterConfig)
	if err != nil {
		return nil, err
	}
	// tmp directory = parent directory plus unique identifier for worker
	wrkr, err := worker.Create(ctx, conv, fmt.Sprintf("%s/worker_%d", conf.DirWork, i))
	if err != nil {
		return nil, err
	}
	return wrkr, nil
}

// i is unique to each worker
type WorkerCallback func(ctx context.Context, conf *Configuration, i int) (worker.Worker, error)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	jobC   chan<- worker.Job
}

// run a manager that listens on a grpc url
func Run(ctx context.Context, binFfmpeg string, binBlender string, args []string) error {
	maxJobs, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	config := &Configuration{
		DirWork:    args[0],
		BinFfmpeg:  binFfmpeg,
		BinBlender: binBlender,
		ListenUrl:  args[1],
		MaxJobs:    maxJobs,
	}

	err = config.Check()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", config.ListenUrl)
	if err != nil {
		return err
	}

	manager, err := Create(ctx, config, DefaultCreateWorker)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	manager.Add(s)

	//TODO: Add go routine to close listener
	go loopClose(ctx, l, s)
	return s.Serve(l)
}

func Create(parentCtx context.Context, conf *Configuration, workerCb WorkerCallback) (Manager, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	// no buffer so program will block if failsc
	// how many jobs do you want in the buffer? for now no buffer
	jobC := make(chan worker.Job, 1)
	// creating worker pool
	n := conf.MaxJobs
	for i := 0; i < n; i++ {
		w, err := workerCb(ctx, conf, i)
		if err != nil {
			cancel()
			return Manager{}, fmt.Errorf("failed to create worker: %s", err)
		} else {
			log.Printf("worker created successfully")
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
	log.Printf("closing loop: %s", err)
}

func loopClose(ctx context.Context, listner net.Listener, s *grpc.Server) {
	<-ctx.Done()
	s.GracefulStop()
	listner.Close()
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
		return result.Reader, result.Err
	case <-m.ctx.Done():
		// The manager's context has been cancelled, return an error
		return nil, m.ctx.Err()
	case <-job.Ctx.Done():
		// The job's context was been cancelled, return an error
		return nil, job.Ctx.Err()

	}
}

// Check if executables (blender and ffmpeg) exist and are executable
func ExecCheck(binFilePath string) error {

	// Check if the file exists
	fi, err := os.Stat(binFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("executable '%s' not found", binFilePath)
		} else {
			panic(fmt.Sprintf("error checking ffmpeg executable: %v", err))
		}
	}
	//Check if the file is executable
	if fi.Mode()&0111 == 0 {
		return fmt.Errorf("executable '%s' is not executable", binFilePath)
	}
	return nil

}
