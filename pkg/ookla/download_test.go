package ookla

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadOnReadWriter(t *testing.T) {
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
