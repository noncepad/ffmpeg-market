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
// args: 0=url string, 1=fileIn string, 2=DirOut string, 3=ext []string
func Run(parentCtx context.Context, args []string) error {

	go func() {
		<-parentCtx.Done()
		log.Print("parent ctx done")
	}()
	// Check if the output directory exists
	if _, err := os.Stat(args[2]); os.IsNotExist(err) {
		return fmt.Errorf("Run - 1: output directory does not exist, err: %s", err)
	}

	// Create a client
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	conn, err := grpc.DialContext(ctx, args[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client := Create(ctx, conn)
	blenderIn := args[1]
	dir := args[2]
	extList := args[3:]

	err = client.ProcessRequest(ctx, blenderIn, dir, extList)
	if err != nil {
		return err
	}

	return err
}
