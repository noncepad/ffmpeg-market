package converter_test

import (
	"os"
	"os/exec"
	"testing"
)

func TestFfmpegConverter(t *testing.T) {
	// Input and output file paths
	inputFilePath := "./files/input.mp4"
	outputFilePath := "./files/output.gif"

	// FFmpeg command
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-f", "gif", "pipe:1 >", outputFilePath)

	// Execute FFmpeg command
	if err := cmd.Run(); err != nil {
		t.Errorf("error converting video to gif: %v", err)
	}

	// Check if the output file exists
	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Errorf("output file %s does not exist", outputFilePath)
	}

}
