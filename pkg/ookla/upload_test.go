package ookla

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRW struct {
	readFn  func(p []byte) (n int, err error)
	writeFn func(p []byte) (n int, err error)
}

func (t testRW) Read(p []byte) (n int, err error) {
	if t.readFn == nil {
		return 0, io.EOF
	}
	return t.readFn(p)
}

func (t testRW) Write(p []byte) (n int, err error) {
	if t.writeFn == nil {
		return 0, io.ErrClosedPipe
	}
	return t.writeFn(p)
}

func TestUploadOnReadWriter(t *testing.T) {
	firstTime := true
	rw := testRW{
		readFn: func(p []byte) (n int, err error) {
			return 0, nil
		},
		writeFn: func(p []byte) (n int, err error) {
			if firstTime {
				firstTime = false
				assert.Equal(t, hiMsg, string(p))
			}
			return len(p), nil
		},
	}

	chunkCh := make(chan int, 1)
	ctx, cancelFn := context.WithCancel(context.Background())
	cancelFn()
	testUploadOnReadWriter(ctx, chunkCh, rw)
}
