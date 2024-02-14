package client

import (
	"context"
	"log"
)

type Log struct {
	Out string
}

func loopLog(ctx context.Context, logC <-chan Log) {
	doneC := ctx.Done()
out:
	for {
		select {
		case <-doneC:
			break out
		case x := <-logC:
			log.Printf("server log: %s", x)
		}
	}
}
