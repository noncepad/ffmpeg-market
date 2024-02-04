package converter_test

import (
	"context"
	"io"
	"os"
	"os/exec"
	"testing"

	pkgcvr "gitlab.noncepad.com/naomiyoko/ffmpeg-market/converter"
)

func TestFFmpegConversion(t *testing.T) {
	// Define input and output file paths
	inputFile := "../files/solpipe.mp4"
	outputFile := "../files/solpipe2.gif"

	// Construct FFmpeg command
	cmd := exec.Command("ffmpeg", "-i", inputFile, outputFile)

	// Run FFmpeg command
	err := cmd.Run()

	// Check if FFmpeg command produced any errors
	if err != nil {
		t.Fatalf("FFmpeg command failed: %v", err)
	}

	// Check if the output file exists after conversion
	_, err = os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Output file does not exist: %v", err)
	}
}

func TestConvert(t *testing.T) {

	// Define input and output file paths
	// inputfile must exist, outputfile must not exist
	inputFile := "../files/solpipe.mp4"
	ext := "gif"
	outputFile := "../files/solpipe." + ext

	// Create a simple converter instance
	ctx := context.Background()
	// make config a reference
	config := &pkgcvr.Configuration{}
	// create a converter with CreateConverter function call
	converter, err := pkgcvr.CreateSimpleConverter(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create converter: %v", err)
	}

	// Call the Convert function
	outputReader, err2 := converter.Convert(ctx, inputFile, ext)
	if err2 != nil {
		t.Fatalf("Convert function returned an error: %v", err)
	}
	// you want to copy to disk, not memory, so we do not use ReadAll (but you are on the right track)
	//_, err = io.ReadAll(outputReader)
	f, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	// the secret trick is here:
	_, err = io.Copy(f, outputReader)
	if err != io.EOF {
		err = nil
	}
	if err != nil {
		t.Fatal(err)
	}

	// Check if the output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file does not exist: %v", err)
	}
}
