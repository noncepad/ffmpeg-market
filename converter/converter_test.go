package converter_test

import (
	"context"
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
	outputFile := "../files/solpipe2.gif"

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
	err = converter.Convert(ctx, inputFile, outputFile)
	if err != nil {
		t.Errorf("Convert function returned an error: %v", err)
	}

	// Check if the output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file does not exist: %v", err)
	}
}
