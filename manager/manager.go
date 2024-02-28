package manager

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/noncepad/ffmpeg-market/converter"
	pbf "github.com/noncepad/ffmpeg-market/proto/ffmpeg"
	"github.com/noncepad/ffmpeg-market/worker"
	poolmgr "github.com/noncepad/worker-pool/manager"
	pool "github.com/noncepad/worker-pool/pool"
	"google.golang.org/grpc"
)

// Settings required for the manager to operate.
type Configuration struct {
	DirWork    string // Directory for work files/temporary files.
	BinFfmpeg  string // Binary location for ffmpeg.
	BinBlender string // Binary location for Blender.
	ListenUrl  string // URL on which to listen for incoming gRPC connections.
	MaxJobs    int    // Maximum number of concurrent jobs/workers.
}

// Returns a copy of the Configuration object.
func (conf *Configuration) Copy() *Configuration {
	x := *conf
	return &x
}

// Checks if configuration properties are set and correctly
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

// Creates instance of worker
func DefaultCreateWorker(ctx context.Context, conf *Configuration, i int) (pool.Worker[worker.Request, worker.Result], error) {
	converterConfig := &converter.Configuration{
		BinFfmpeg:  conf.BinFfmpeg,
		BinBlender: conf.BinBlender,
	}
	// Creates a new converter based on the passed context and converter configuration.
	conv, err := converter.CreateSimpleConverter(ctx, converterConfig)
	if err != nil {
		return nil, err
	}
	// Sets up a tmp directory = parent directory plus unique identifier for worker
	wrkr, err := worker.Create(ctx, conv, fmt.Sprintf("%s/worker_%d", conf.DirWork, i))
	if err != nil {
		return nil, err
	}
	return wrkr, nil
}

// Sets up the gRPC manager and starts listening for incoming connections.
func Run(ctx context.Context, binFfmpeg string, binBlender string, args []string) error {
	// Parse the maximum number of jobs from the command line arguments.
	maxJobs, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	// Create a new manager for handling jobs with the parsed number of maxJobs.
	manager, err := poolmgr.Create[worker.Request, worker.Result](ctx, maxJobs)
	if err != nil {
		return err
	}
	// Set up configuration based on passed arguments and binaries.
	config := &Configuration{
		DirWork:    args[0],
		BinFfmpeg:  binFfmpeg,
		BinBlender: binBlender,
		ListenUrl:  args[1],
		MaxJobs:    maxJobs,
	}
	// Check for correct configuration parameters
	err = config.Check()
	if err != nil {
		return err
	}
	// For each job, instantiate a worker and add it to the manager.
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
	// Start listening on TCP based on the configured ListenUrl.
	l, err := net.Listen("tcp", config.ListenUrl)
	if err != nil {
		return err
	}

	s := grpc.NewServer() // Create a new gRPC server instance.
	// Register the JobManagerServer service with the manager
	pbf.RegisterJobManagerServer(s, external{manager: manager})
	go loopClose(ctx, l, s) //go routine to close listener

	return s.Serve(l) // Serve the gRPC server on the established listener.

}

// Wait for cancellation and then stops the server and closes the listener.
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
