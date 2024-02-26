package manager

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	poolmgr "github.com/noncepad/worker-pool/manager"
	pool "github.com/noncepad/worker-pool/pool"
	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/converter"
	pbf "gitlab.noncepad.com/naomiyoko/ffmpeg-market/proto/ffmpeg"
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

func (conf *Configuration) Copy() *Configuration {
	x := *conf
	return &x
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
func DefaultCreateWorker(ctx context.Context, conf *Configuration, i int) (pool.Worker[worker.Request, worker.Result], error) {
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

// run a manager that listens on a grpc url
func Run(ctx context.Context, binFfmpeg string, binBlender string, args []string) error {

	maxJobs, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}

	manager, err := poolmgr.Create[worker.Request, worker.Result](ctx, maxJobs)
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

	for i := 0; i < maxJobs; i++ {

		w, err := DefaultCreateWorker(ctx, config, i)
		if err != nil {
			manager.Close()
			return err
		}
		err = manager.Add(w)
		if err != nil {
			manager.Close()
			return err
		}
	}

	l, err := net.Listen("tcp", config.ListenUrl)
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	pbf.RegisterJobManagerServer(s, external{manager: manager})
	//TODO: Add go routine to close listener
	go loopClose(ctx, l, s)
	return s.Serve(l)
}

func loopClose(ctx context.Context, listner net.Listener, s *grpc.Server) {
	<-ctx.Done()
	s.GracefulStop()
	listner.Close()
}

type clientLogger interface {
	Send(string)
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
