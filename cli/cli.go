package cli

import (
	"fmt"
	"io"
)

type Executor interface {
	Execute(args []string, stdin io.Reader, stdout, stderr io.Writer) (exitCode int)
}

func MustPrintf(w io.Writer, format string, args ...any) {
	if _, err := fmt.Fprintf(w, format, args...); err != nil {
		panic(err)
	}
}

func PanicIfErrorf(err error, format string, args ...any) {
	if err != nil {
		panic(fmt.Errorf(format+": %w", append(args, err)...))
	}
}
