package worker_test

import (
	"context"
	"testing"

	pkgcvr "gitlab.noncepad.com/naomiyoko/ffmpeg-market/converter"
	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/worker"
	pkgwkr "gitlab.noncepad.com/naomiyoko/ffmpeg-market/worker"
)

func TestRun(t *testing.T) {
	// Define a temporary directory for testing
	tmpDir := "/tmp"

	// Create a simple converter instance
	ctx := context.Background()
	config := &pkgcvr.Configuration{}
	converter, err := pkgcvr.CreateSimpleConverter(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create converter: %v", err)
	}

	// Create a worker
	worker, err := worker.Create(ctx, converter, tmpDir)
	if err != nil {
		t.Fatalf("Failed to create worker: %v", err)
	}

	// Define a job
	job := pkgwkr.Job{
		Ctx:     ctx,
		Blender: "solpop.blend",
		Out:     []string{"mp4", "mkv", "gif"},
		ResultC: make(chan<- pkgwkr.Result),
	}

	// Run the job
	result := worker.Run(job)

	// Check for errors in the result
	if result.Err != nil {
		t.Fatalf("Run function returned an error: %v", result.Err)
	}

	// Check if the readers are not nil
	for _, reader := range result.Reader {
		if reader == nil {
			t.Fatal("Reader is nil")
		}
	}

	// Close the readers after use
	for _, reader := range result.Reader {
		reader.Close()
	}
}
