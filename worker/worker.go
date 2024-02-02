package worker

import (
	"context"
	"io"
)

type Job struct {
	Ctx     context.Context
	Blender string   // filepath
	Out     []string // list of file extensions
	ResultC chan<- Result
}

type Result struct {
	Reader []io.ReadCloser // file handle
	Err    error
}

// err: if render fails -> no file handle for conversion
type Worker interface {
	Run(Job) Result
	// this is ther workers version of ctx.Done()
	CloseSignal() <-chan error
}
