package ookla

import (
	"context"
	"fmt"
	"io"
)

const (
	uploadTpl       = "UPLOAD %d 0\n"
	uploadChunkSize = 64 * 1024
)

// testUpload tests upload from server and sends sizes of uploaded chunks into the channel
func (s *Server) testUpload(ctx context.Context, chunksCh chan<- int) error {
	conn, err := s.dial(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", s.tcpAddr.String(), err)
	}
	defer conn.Close()

	return testUploadOnReadWriter(ctx, chunksCh, conn)
}

// created for testing purpose
func testUploadOnReadWriter(ctx context.Context, chunksCh chan<- int, rw io.ReadWriter) error {
	rw.Write([]byte(hiMsg))
	msg := make([]byte, 1024)
	rw.Read(msg)

	requested := uploadChunkSize

	// according to what i found ookla's servers understand such messages
	// to prepare server for uploading to it
	uploadMsg := []byte(fmt.Sprintf(uploadTpl, requested))
	msg = make([]byte, requested-len(uploadMsg))

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		_, err := rw.Write(uploadMsg)
		if err != nil {
			return fmt.Errorf("failed to upload message: %w", err)
		}
		_, err = rw.Write(msg)
		if err != nil {
			return fmt.Errorf("failed to upload data: %w", err)
		}

		_, err = rw.Read(msg)
		if err != nil {
			return fmt.Errorf("failed to download after upload data: %w", err)
		}

		chunksCh <- requested
	}
}
