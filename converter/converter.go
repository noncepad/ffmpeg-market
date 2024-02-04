package converter

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Converter interface {
	// takes an mp4 inputfile and converts it to a gif output file
	Convert(ctx context.Context, inputFile string, ext string) (io.ReadCloser, error)
	Render(ctx context.Context, inputFile, outputFile string) error
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
func (sc *simpleConverter) Convert(ctx context.Context, inputFile string, ext string) (io.ReadCloser, error) {
	// use create to make sure we are not overwriting a file
	// the command will fail if the output file already exists
	saveHandle, stdout := io.Pipe()

	// OLD: ffmpeg -i ./files/solpipe.mkv ./files/solpipe2.gif
	// new: ffmpeg -i ./files/solpipe.mp4 -f gif pipe:1 > ./files/solpipe5.gif
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", inputFile, "-f", ext, "pipe:1")

	// redirect stdout to the save handle
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		err = fmt.Errorf("ffmpeg command failed: %v", err)
		return nil, err
	}

	fmt.Println("Conversion successful")
	return saveHandle, cmd.Wait()
}
func (sc *simpleConverter) Render(ctx context.Context, inputFile, outputFile string) error {
	// use create to make sure we are not overwriting a file
	// the command will fail if the output file already exists
	saveFileHandle, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	// blender -b ./files/solpop.blend -o ./files/solpop.avi -F AVIJPEG -x 1 -f 1 -a
	cmd := exec.CommandContext(ctx, "blender", "--b", inputFile, "-out", outputFile, "-F", "AVIJPEG", "-x", "1", "-f", "1", "-a")

	// redirect stdout to the save handle
	cmd.Stdout = saveFileHandle
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("blender command failed: %v", err)
		return err
	}

	fmt.Println("Rendering successful")
	return cmd.Wait()
}
