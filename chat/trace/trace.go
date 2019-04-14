package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

type nilTracer struct {
}

func (t *tracer) Trace(args ...interface{})  {
	t.out.Write([]byte(fmt.Sprint(args...)))
	t.out.Write([]byte("\n"))
}
func (t *nilTracer) Trace(args ...interface{})  {
	// noop
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
func Noop() Tracer {
	return &nilTracer{}
}