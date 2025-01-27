package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StringWriteCloser struct {
	bytes.Buffer
}

func (swc *StringWriteCloser) Close() error {
	// No resources to release, so just return nil
	return nil
}

func (swc *StringWriteCloser) Write(p []byte) (n int, err error) {
	return swc.Buffer.Write(p)
}

func (swc *StringWriteCloser) String() string {
	return swc.Buffer.String()
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		want    *Log
		options []LoggingOption
	}{
		{
			name: "default",
			want: &Log{
				debug:  false,
				output: os.Stderr,
			},
			options: nil,
		},
		{
			name: "with debug",
			want: &Log{
				debug:  true,
				output: os.Stderr,
			},
			options: []LoggingOption{WithDebug(true)},
		},
		{
			name: "with output",
			want: &Log{
				debug:  false,
				output: os.Stdout,
			},
			options: []LoggingOption{WithOutput(os.Stdout)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.options...)
			assert.Equal(t, l, tt.want)
		})
	}
}

func TestDebug(t *testing.T) {
	tests := []struct {
		name      string
		debug     bool
		message   string
		additions []interface{}
		expected  string
	}{
		{
			name:     "debug off",
			debug:    false,
			message:  "this is a debug message",
			expected: "",
		},
		{
			name:     "debug on",
			debug:    true,
			message:  "this is a debug message",
			expected: "[DEBUG] this is a debug message\n",
		},
		{
			name:      "debug on with additional arguments",
			debug:     true,
			message:   "this is a debug message with %d arguments",
			additions: []interface{}{3},
			expected:  "[DEBUG] this is a debug message with 3 arguments\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			swc := &StringWriteCloser{}

			options := []LoggingOption{
				WithOutput(swc),
				WithDebug(tt.debug),
			}

			l := New(options...)
			if len(tt.additions) == 0 {
				l.Debug(tt.message)
			} else {
				l.Debugf(tt.message, tt.additions...)
			}
			assert.True(t, strings.HasSuffix(swc.String(), tt.expected))

		})
	}
}
