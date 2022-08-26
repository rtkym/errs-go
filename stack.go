package errs

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
)

func getStackTrace() StackTrace {
	const (
		maxDepth    = 64
		callerDepth = 4
	)

	pcs := make([]uintptr, maxDepth)
	n := runtime.Callers(callerDepth, pcs)

	resp := make([]StackFrame, n)

	for i := 0; i < n; i++ {
		resp[i] = frame(pcs[i]).toStackFrame()
	}

	return resp
}

type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) - 1 }

func (f frame) toStackFrame() StackFrame {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return StackFrame{Func: "unknown", File: "unknown", Line: 0}
	}

	file, line := fn.FileLine(f.pc())

	return StackFrame{
		Func: fn.Name(),
		File: file,
		Line: line,
	}
}

type StackTrace []StackFrame

func (s StackTrace) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case state.Flag('+'):
			for _, f := range s {
				fmt.Fprintf(state, "\n%+v", f)
			}
		default:
		}
	default:
	}
}

type StackFrame struct {
	Func string `json:"func"`
	File string `json:"file"`
	Line int    `json:"line"`
}

func (f StackFrame) Format(state fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case state.Flag('+'):
			_, _ = io.WriteString(state, f.Func)
			_, _ = io.WriteString(state, "\n\t")
			_, _ = io.WriteString(state, f.File)
		default:
			_, _ = io.WriteString(state, path.Base(f.File))
		}
	case 'd':
		_, _ = io.WriteString(state, strconv.Itoa(f.Line))
	case 'v':
		f.Format(state, 's')
		_, _ = io.WriteString(state, ":")
		f.Format(state, 'd')
	}
}
