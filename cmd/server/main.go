package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/noncepad/ffmpeg-market/manager"
)

// if you dont set the filepaths as enviromnetal variables this func will lookup executable in the PATH
// ./server <myworkingdirectory> <listenurl for grpc> <MaxJobs>
func main() {
	// Set up channel to listen for interrupt or termination signals
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGTERM, syscall.SIGINT)
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go loopSignal(ctx, cancel, signalC) // Start a goroutine that listens for OS signals and cancels the context when received.

	var err error

	// Check for the correct number of Args
	if len(os.Args) != 4 {
		panic("main - 1: not the correct number of arguments")
	}

	// which ffmpeg (find ffmpeg executable)
	whichFfmpeg, present := os.LookupEnv("BIN_FFMPEG")
	if !present {
		whichFfmpeg, err = exec.LookPath("ffmpeg")
		if err != nil {
			panic("main - 2: ffmpeg executable not found in PATH")
		}
	}

	// which blender (find blender executable)
	whichBlender, present := os.LookupEnv("BIN_BLENDER")
	if !present {
		whichBlender, err = exec.LookPath("blender")
		if err != nil {
			panic("main - 3: blender executable not found in PATH")
		}
	}

	// Run the manager with the given paths and command line arguments.
	log.Println("main - 4: Sufficient arguments, running manager")
	err = manager.Run(ctx, whichFfmpeg, whichBlender, os.Args[1:])
	if err != nil {
		panic(fmt.Sprintf("main - 8: Error running manager: %v", err))
	}

}

func loopSignal(ctx context.Context, cancel context.CancelFunc, signalC <-chan os.Signal) {
	defer cancel()
	doneC := ctx.Done()
	select {
	case <-doneC:
	case s := <-signalC:
		// Received an OS signal; cancel the context and write to stderr.
		os.Stderr.WriteString(fmt.Sprintf("%s\n", s.String()))
		log.Printf("loopSignal: Received signal: %s", s.String())
	}
}
