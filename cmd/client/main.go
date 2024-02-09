package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/client"
)

// ./client <listenurl> <filepathIn> <DirectoryOut> <extesion1> <extesion2> .... <extesionN>
func main() {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go loopSignal(ctx, cancel, signalC)
	// send blender to server (run)

	if len(os.Args) < 5 {
		panic("main - 1: insufficient arguments")
	}

	log.Println("main - 2: Sufficient arguments, running client")

	// Call clients Run function
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
	case s := <-signalC:
		log.Printf("loopSignal: Received signal: %s\n", s.String())
	}
}
