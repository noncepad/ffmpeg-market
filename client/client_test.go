package client_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/noncepad/ffmpeg-market/client"
	"github.com/noncepad/ffmpeg-market/manager"
)

// server command: server /tmp/myworkdir/ localhost:30051 3

// client command: client localhost:30051 ../files/solpop.blend ../files mkv mp4 gif

// TestFfmpeg is a test function for testing the ffmpeg functionality
func TestFfmpeg(t *testing.T) {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure we cancel the context when this function returns

	// Clean up output files before running tests
	for _, v := range []string{"gif", "mp4", "mpeg", "ogv"} {
		os.Remove("../files/out." + v) // Remove file if it exists
	}

	// URL of the server
	url := "localhost:30051"
	// Channel for receiving errors from the concurrently running manager
	errorC := make(chan error, 1)

	// Remove any existing directory and recreate it
	os.RemoveAll("/tmp/myworkdir")
	err := os.Mkdir("/tmp/myworkdir", 0755)
	if err != nil {
		t.Fatal(err) // If there's an error creating the directory, stop the test
	}
	// Cleanup function to remove the temporary work directory after the test completes
	t.Cleanup(func() { os.RemoveAll("/tmp/myworkdir") })

	// Start the loopRunManager function in a new goroutine
	go loopRunManager(ctx, errorC, url, cancel)
	time.Sleep(10 * time.Second) // Wait for the manager to start and settle

	// Call the client.Run function, which will attempt to connect to the server and run ffmpeg processes
	err = client.Run(ctx, []string{url, "../files/solpop.blend", "../files", "ogv", "mp4", "gif"})
	if err != nil {
		t.Fatal(err) // If there's an error, stop the test
	}

}

// loopRunManager is a helper function to initiate the manager and send any errors back through a channel
func loopRunManager(ctx context.Context, errorC chan<- error, url string, cancel context.CancelFunc) {
	defer cancel() // Cancel context on exit to ensure cleanup

	// Call manager.Run, passing the context, program paths, and configuration options, then send any errors to the errorC channel
	errorC <- manager.Run(ctx, "/usr/bin/ffmpeg", "/usr/bin/blender", []string{"/tmp/myworkdir", url, "3"})

	// Log a message when the manager stops running
	log.Printf("exiting manager")
}
