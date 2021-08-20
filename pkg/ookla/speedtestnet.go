package ookla

import (
	"context"
	"fmt"
	"time"

	"github.com/go-speedtest/pkg/speed"
)

type provider string

// Provider is the constant for Ookla speed test
const Provider provider = "speedtest.net"

// Name is Ookla's name
const Name = "ookla"

func (t provider) Test(ctx context.Context, resultCh chan<- speed.Measures) error {
	config, err := GetConfiguration()
	if err != nil {
		return err
	}

	servers, err := GetServers(config)
	if err != nil {
		return err
	}

	measures := speed.Measures{}

	measures.Download = "*"
	resultCh <- measures

	downloadCtx, cancelFn := context.WithTimeout(ctx, time.Duration(config.Download.Timeout)*time.Second)
	defer cancelFn()

	download := transferStrategy{
		servers:  servers,
		resultCh: resultCh,
		updateMeasures: func(measures speed.Measures, bytesPerSecond float64) speed.Measures {
			measures.Download, measures.DownloadUnit = speed.Format(bytesPerSecond)
			return measures
		},
		runTransfer: func(ctx context.Context, server Server, chunkCh chan<- int) error {
			return server.testDownload(ctx, chunkCh)
		},
	}
	measures, err = download.testTransfer(downloadCtx, measures)
	if err != nil {
		return fmt.Errorf("download test failed: %w", err)
	}

	measures.Upload = "*"
	resultCh <- measures

	uploadCtx, cancelFn := context.WithTimeout(ctx, time.Duration(config.Upload.Timeout)*time.Second)
	defer cancelFn()

	upload := transferStrategy{
		servers:  servers,
		resultCh: resultCh,
		updateMeasures: func(measures speed.Measures, bytesPerSecond float64) speed.Measures {
			measures.Upload, measures.UploadUnit = speed.Format(bytesPerSecond)
			return measures
		},
		runTransfer: func(ctx context.Context, server Server, chunkCh chan<- int) error {
			return server.testUpload(ctx, chunkCh)
		},
	}
	measures, err = upload.testTransfer(uploadCtx, measures)
	if err != nil {
		return fmt.Errorf("upload test failed: %w", err)
	}

	close(resultCh)

	return nil
}

func (t provider) String() string {
	return string(t)
}
