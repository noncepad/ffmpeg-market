package converter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type Converter interface {
	// takes an mp4 inputfile and converts it to a gif output file
	Convert(ctx context.Context, inputFile, outputFile string) error
}

type simpleConverter struct {
	ctx    context.Context
	config *Configuration
}

type Configuration struct {
}

func CreateSimpleConverter(ctx context.Context, config *Configuration) (Converter, error) {
	e1 := new(simpleConverter)
	e1.ctx = ctx
	e1.config = config
	return e1, nil
}

// we need to put a context variable here
// the sc.ctx is for the server as a whole.
// ctx is for the user request.
func (sc *simpleConverter) Convert(ctx context.Context, inputFile, outputFile string) error {
	// use create to make sure we are not overwriting a file
	// the command will cail if the output file already exists
	saveFileHandle, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	// OLD: ffmpeg -i ./files/solpipe.mkv ./files/solpipe2.gif
	// new: ffmpeg -i ./files/solpipe.mp4 -f gif pipe:1 > ./files/solpipe5.gif
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", inputFile, "-f", "gif", "pipe:1")

	// redirect stdout to the save handle
	cmd.Stdout = saveFileHandle
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("ffmpeg command failed: %v", err)
		return err
	}

	fmt.Println("Conversion successful")
	return cmd.Wait()
}
