package client

import (
	"context"
	"log"
)

type Log struct {
	Out string // Out is a field that stores the log message.
}

func loopLog(ctx context.Context, logC <-chan Log) {
	doneC := ctx.Done() // doneC will receive a signal when the context is cancelled.
out:
	for {
		select {
		case <-doneC:
			// If the context has been cancelled, break out of the loop.
			break out
		case x := <-logC:
			// Receive a log entry from the channel and print it.
			log.Printf("server log: %s", x)
		}
	}
}
