package client_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/client"
	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/manager"
)

// server command: server /tmp/myworkdir/ localhost:30051 3

// client command: client localhost:30051 ../files/solpop.blend ../files mkv mp4 gif

func TestFfmpeg(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, v := range []string{"gif", "mp4", "mpeg"} {
		os.Remove("../files/out." + v)
	}

	url := "localhost:30051"
	errorC := make(chan error, 1)
	os.RemoveAll("/tmp/myworkdir")
	err := os.Mkdir("/tmp/myworkdir", 0755)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll("/tmp/myworkdir") })
	go loopRunManager(ctx, errorC, url, cancel)
	time.Sleep(10 * time.Second)
	err = client.Run(ctx, []string{url, "../files/solpop.blend", "../files", "mp4", "mpeg"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("why am i here??")
}

func loopRunManager(ctx context.Context, errorC chan<- error, url string, cancel context.CancelFunc) {
	defer cancel()
	errorC <- manager.Run(ctx, "/usr/bin/ffmpeg", "/usr/bin/blender", []string{"/tmp/myworkdir", url, "3"})
	log.Printf("exiting manager")
}
