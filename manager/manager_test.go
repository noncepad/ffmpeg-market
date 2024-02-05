package manager_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/converter"
	pkgmngr "gitlab.noncepad.com/naomiyoko/ffmpeg-market/manager"
	"gitlab.noncepad.com/naomiyoko/ffmpeg-market/worker"
)

func TestSendJob(t *testing.T) {
	ctx := context.Background()
	config := &pkgmngr.Configuration{
		DirWork:    "/tmp",
		BinFfmpeg:  "/usr/bin/ffmpeg",
		BinBlender: "/usr/bin/blender",
		ListenUrl:  "blah",
		MaxJobs:    1,
	}

	err := config.Check()
	if err != nil {
		t.Fatal(err)
	}

	manager, err := pkgmngr.Create(ctx, config, createWorker)
	if err != nil {
		t.Fatal(err)
	}

	blender := "../files/solpop.blend"
	out := []string{"mp4", "mkv", "gif"}
	var responseList []io.ReadCloser

	expectedNumResponses := len(out)

	responseList, err = manager.SendJob(ctx, blender, out)
	if err != nil {
		t.Fatal(err)
	}
	if responseList == nil {
		t.Fatalf("Expected non-nil responseList, got nil: %s", err)
	}
	if len(responseList) != expectedNumResponses {
		t.Errorf("Expected %d responses in responseList, got %d", expectedNumResponses, len(responseList))
	}

}

func createWorker(ctx context.Context, conf *pkgmngr.Configuration, i int) (worker.Worker, error) {
	converterConfig := &converter.Configuration{
		BinFfmpeg:  conf.BinFfmpeg,
		BinBlender: conf.BinBlender,
	}

	conv, err := converter.CreateSimpleConverter(ctx, converterConfig)
	if err != nil {
		return nil, err
	}
	// tmp directory = parent directory plus unique identifier for worker
	wrkr, err := worker.Create(ctx, conv, fmt.Sprintf("%s/worker_%d", conf.DirWork, i))
	if err != nil {
		return nil, err
	}
	return wrkr, nil
}
