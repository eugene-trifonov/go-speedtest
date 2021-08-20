package ookla

import (
	"context"
	"fmt"
	"io"
)

const (
	downloadTpl = "DOWNLOAD %d\n"
)

// testDownload tests download from server and sends sizes of downloaded chunks into the channel
func (s *Server) testDownload(ctx context.Context, chunksCh chan<- int) error {
	conn, err := s.dial(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", s.tcpAddr.String(), err)
	}
	defer conn.Close()
	return testDownloadOntoReadWriter(ctx, chunksCh, conn)
}

// created for testing purpose
func testDownloadOntoReadWriter(ctx context.Context, chunksCh chan<- int, rw io.ReadWriter) error {
	rw.Write([]byte(hiMsg))
	hello := make([]byte, 1024)
	rw.Read(hello)

	var downloaded int

	tmp := make([]byte, 2*1024)

	requested := 64 * 1024

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// according to what I found ookla's servers understand such messages
		// to prepare server for downloading from it
		rw.Write([]byte(fmt.Sprintf(downloadTpl, requested)))

		for downloaded < requested {
			n, err := rw.Read(tmp)
			if err != nil {
				switch err {
				case io.EOF:
				default:
					return fmt.Errorf("failed to download data: %w", err)
				}
				break
			}
			downloaded += n
		}

		chunksCh <- downloaded

		downloaded = 0
		requested *= 2
	}
}
