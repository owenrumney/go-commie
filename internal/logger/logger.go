package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Log struct {
	debug  bool
	output io.WriteCloser
}

type LoggingOption func(*Log)

func WithDebug(useDebug bool) LoggingOption {
	return func(l *Log) {
		l.debug = useDebug
	}
}

func WithOutput(w io.WriteCloser) LoggingOption {
	return func(l *Log) {
		l.output = w
	}
}

func New(options ...LoggingOption) *Log {
	l := &Log{
		debug:  false,
		output: os.Stderr,
	}

	for _, opt := range options {
		opt(l)
	}

	return l
}

func (l *Log) Debug(msg string) {
	l.Debugf("%s", msg)
}

func (l *Log) Debugf(msg string, args ...interface{}) {
	if l.debug {
		msg = fmt.Sprintf("%s [DEBUG] %s\n", time.Now().Format(time.RFC3339Nano), msg)
		fmt.Fprintf(l.output, msg, args...)
	}
}

func (l *Log) Fatal(err error) {
	fmt.Fprintf(l.output, "error: %s\n", err)
	os.Exit(1)
}
