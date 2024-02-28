package client

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// host:port
// Run starts the client operation with the given arguments.
// args: 0=url string, 1=fileIn string, 2=DirOut string, 3=ext []string
func Run(parentCtx context.Context, args []string) error {
	// Start a goroutine that will print a message once the parent context is done.
	go func() {
		<-parentCtx.Done()
		log.Print("parent ctx done")
	}()
	// Check if the output directory exists
	if _, err := os.Stat(args[2]); os.IsNotExist(err) {
		return fmt.Errorf("Run - 1: output directory does not exist, err: %s", err)
	}

	// Create a new context with cancel functionality from the parent context
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	// Dial the gRPC server using provided URL and insecure credentials as no TLS config is assumed
	conn, err := grpc.DialContext(ctx, args[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client := Create(ctx, conn) // returns a client interface with ProcessRequest method.
	blenderIn := args[1]
	dir := args[2]
	extList := args[3:]

	// Call ProcessRequest on the client with the necessary arguments.
	err = client.ProcessRequest(ctx, blenderIn, dir, extList)
	if err != nil {
		return err
	}

	return err
}
