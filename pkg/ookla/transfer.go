package ookla

import (
	"context"
	"time"

	"github.com/go-speedtest/pkg/speed"
	"golang.org/x/sync/errgroup"
)

const hiMsg = "Hi\n"

// transferStrategy avoids a duplication of download and upload testings
// it uses two functions to differentiate download nad upload strategies
type transferStrategy struct {
	servers  []Server
	resultCh chan<- speed.Measures

	updateMeasures func(measures speed.Measures, bytesPerSecond float64) speed.Measures
	runTransfer    func(ctx context.Context, server Server, chunksCh chan<- int) error
}

// testTransfer is a main transfer strategy function
// it is similar for download nad upload
func (s transferStrategy) testTransfer(ctx context.Context, measures speed.Measures) (speed.Measures, error) {
	chunksCh := make(chan int, len(s.servers))

	startTime := time.Now()

	gr, ctx := errgroup.WithContext(ctx)
	for i := range s.servers {
		server := s.servers[i]
		gr.Go(func() error {
			return s.runTransfer(ctx, server, chunksCh)
		})
	}

	var total int
	for {
		select {
		case <-ctx.Done():
			// at this point either context is Canceled or Timeouted, error group is also should be over
			// we need to get an error from the group because it can be a reason why context is cancelled
			return measures, gr.Wait()
		case chunk := <-chunksCh:
			total += chunk
			bytesPerSecond := float64(total) / float64(time.Now().Sub(startTime).Seconds())
			measures = s.updateMeasures(measures, bytesPerSecond)
			s.resultCh <- measures
		}
	}
}
