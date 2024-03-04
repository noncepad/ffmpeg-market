package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/noncepad/ffmpeg-market/client"
)

// ./client <serverURL> <filepathIn> <DirectoryOut> <extesion1> <extesion2> .... <extesionN>
func main() {
	// Set up channel to listen for interrupt or termination signals
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGTERM, syscall.SIGINT)
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go loopSignal(ctx, cancel, signalC) // Start a goroutine that listens for OS signals and cancels the context when received.

	// Check for minimum number of command line arguments
	if len(os.Args) < 5 {
		panic("main - 1: insufficient arguments")
	}

	log.Println("main - 2: Sufficient arguments, running client")

	// Call client's Run function with the command line arguments excluding the program name.
	err := client.Run(ctx, os.Args[1:])
	if err != nil {
		panic(fmt.Sprintf("main - 3: Error running client %v", err))
	}

	log.Println("main - 4: Client execution completed successfully")
}

func loopSignal(ctx context.Context, cancel context.CancelFunc, signalC <-chan os.Signal) {
	defer cancel()
	doneC := ctx.Done()
	select {
	case <-doneC:
	case s := <-signalC: // Received an OS signal.s
		log.Printf("loopSignal: Received signal: %s\n", s.String())
	}
}
