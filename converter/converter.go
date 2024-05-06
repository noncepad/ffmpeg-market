package converter

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Define methods for converting files using ffmpeg and blender executables
type Converter interface {
	Convert(ctx context.Context, inputFile string, outFp string) error
	Render(ctx context.Context, inputFile, outputFile string) error
}

// Implements the Converter interface
type simpleConverter struct {
	ctx    context.Context
	config *Configuration
}

// Hold paths to external tools needed by simpleConverter.
type Configuration struct {
	BinFfmpeg  string
	BinBlender string
}

// Creates an instance of simpleConverter
func CreateSimpleConverter(ctx context.Context, config *Configuration) (Converter, error) {
	e1 := new(simpleConverter)
	e1.ctx = ctx
	e1.config = config
	return e1, nil
}

// Executes FFmpeg to convert files from an .avi format to another video format
func (sc *simpleConverter) Convert(ctx context.Context, inputFile string, outFp string) error {
	// ffmpeg -i ./files/solpipe.mp4 -f gif pipe:1 > ./files/solpipe5.gif
	var err error
	cmd := exec.CommandContext(ctx, sc.config.BinFfmpeg, "-i", inputFile, "-f", strings.TrimPrefix(filepath.Ext(outFp), "."), outFp)

	// redirect stdout to the save handle
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Printf("ffmpeg command failed: %v", err)
		err = fmt.Errorf("ffmpeg command failed: %v", err)
		return err
	}

	log.Printf("ffmpeg %s %s successfully started", inputFile, outFp)
	err = cmd.Wait()
	log.Printf("ffmpeg ext %s done: %s", outFp, err)
	if err != nil {
		log.Printf("ffmpeg command failed: %v", err)
		err = fmt.Errorf("ffmpeg command failed: %v", err)
		return err
	}

	return nil
}

func loopCloseWithError(
	cmd *exec.Cmd,
	saveHandle *io.PipeReader,
) {
	err := cmd.Wait()
	if err != nil {
		saveHandle.CloseWithError(err)
	} else {
		saveHandle.Close()
	}
}

// Executes Blender to convert files from .blend format to an avi. format
func (sc *simpleConverter) Render(ctx context.Context, inputFile, outputFile string) error {
	// use create to make sure we are not overwriting a file
	// the command will fail if the output file already exists
	saveFileHandle, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	// blender -b ./files/solpop.blend -o ./files/solpop.avi -F AVIJPEG -x 1 -f 1 -a
	cmdStr := []string{sc.config.BinBlender, "-b", inputFile, "-o", outputFile, "-F", "AVIJPEG", "-x", "1", "-f", "1", "-a"}
	cmd := exec.CommandContext(ctx, sc.config.BinBlender, "-b", inputFile, "-o", outputFile, "-F", "AVIJPEG", "-x", "1", "-f", "1", "-a")

	// redirect stdout to the save handle
	cmd.Stdout = saveFileHandle
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("blender command failed: %v", err)
		log.Print(err)
		return err
	}

	log.Printf("Rendering started: %+v", cmdStr)
	return cmd.Wait()
}
