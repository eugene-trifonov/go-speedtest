package ookla

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-speedtest/pkg/speed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransferStrategy_NoOps(t *testing.T) {
	noOps := transferStrategy{
		servers: []Server{
			{
				Country: "BY",
				ID:      1,
			},
			{
				Country: "RU",
				ID:      2,
			},
		},
		resultCh: make(chan speed.Measures, 2),
		updateMeasures: func(measures speed.Measures, bytesPerSecond float64) speed.Measures {
			return measures
		},
		runTransfer: func(ctx context.Context, server Server, chunksCh chan<- int) error {
			chunksCh <- 100500
			return nil
		},
	}

	// ookla's speed test is based on time, hence we need to restrict the time
	ctx, cancelFn := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancelFn()

	measures := speed.Measures{Download: "123", DownloadUnit: "Gbps", Upload: "321", UploadUnit: "Mbps"}
	resultMeasures, err := noOps.testTransfer(ctx, measures)
	require.NoError(t, err)
	assert.Equal(t, measures, resultMeasures)
}

func TestTransferStrategy_ErrorFromGoRoutine(t *testing.T) {
	expectedErr := errors.New("error!!!")

	noOps := transferStrategy{
		servers: []Server{
			{
				Country: "BY",
				ID:      1,
			},
			{
				Country: "RU",
				ID:      2,
			},
		},
		resultCh: make(chan speed.Measures, 2),
		updateMeasures: func(measures speed.Measures, bytesPerSecond float64) speed.Measures {
			return measures
		},
		runTransfer: func(ctx context.Context, server Server, chunksCh chan<- int) error {
			return expectedErr
		},
	}

	// ookla's speed test is working on time base, hence we need to restrict the time
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFn()

	measures := speed.Measures{Download: "123", DownloadUnit: "Gbps", Upload: "321", UploadUnit: "Mbps"}
	resultMeasures, err := noOps.testTransfer(ctx, measures)
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, measures, resultMeasures)
}

func TestTransferStrategy_UpdateMeasures(t *testing.T) {
	noOps := transferStrategy{
		servers: []Server{
			{
				Country: "BY",
				ID:      1,
			},
			{
				Country: "RU",
				ID:      2,
			},
		},
		resultCh: make(chan speed.Measures, 2),
		updateMeasures: func(measures speed.Measures, bytesPerSecond float64) speed.Measures {
			measures.Download, measures.DownloadUnit = speed.Format(bytesPerSecond)
			return measures
		},
		runTransfer: func(ctx context.Context, server Server, chunksCh chan<- int) error {
			chunksCh <- 10
			return nil
		},
	}

	// ookla's speed test is working on time base, hence we need to restrict the time
	ctx, cancelFn := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFn()

	measures := speed.Measures{Download: "123", DownloadUnit: "Gbps", Upload: "321", UploadUnit: "Mbps"}
	resultMeasures, err := noOps.testTransfer(ctx, measures)
	require.NoError(t, err)
	assert.NotEqual(t, measures, resultMeasures)
}
