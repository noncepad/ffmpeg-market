package converter_test

import (
	"os"
	"os/exec"
	"testing"
)

func TestFFmpegConversion(t *testing.T) {
	// Define input and output file paths
	inputFile := "/home/naomi/work/ffmpeg-market/local-files/solpipe.mkv"
	outputFile := "/home/naomi/work/ffmpeg-market/local-files/solpipe.gif"

	// Construct FFmpeg command
	cmd := exec.Command("ffmpeg", "-i", inputFile, outputFile)

	// Run FFmpeg command
	err := cmd.Run()

	// Check if FFmpeg command produced any errors
	if err != nil {
		t.Errorf("FFmpeg command failed: %v", err)
	}

	// Check if the output file exists after conversion
	_, err = os.Stat(outputFile)
	if err != nil {
		t.Errorf("Output file does not exist: %v", err)
	}
}
