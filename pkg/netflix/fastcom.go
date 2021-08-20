package netflix

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"

	"github.com/go-speedtest/pkg/speed"
)

var (
	errUpdatesFinished = errors.New("updates are finished")
)

type provider string

const Provider provider = "fast.com"
const Name = "netflix"

// Run starts to collect speed test results
// results are sent back to channel
func (t provider) Test(ctx context.Context, ch chan<- speed.Measures) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err := testSpeed(ctx, ch)
	if err != nil {
		return err
	}
	close(ch)

	return nil
}

func (t provider) String() string {
	return string(t)
}

func testSpeed(ctx context.Context, ch chan<- speed.Measures) error {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	cmds := []chromedp.Action{
		emulation.SetUserAgentOverride(`chromedp/chromedp v0.6.10`),
		chromedp.Navigate(`https://fast.com`),
		chromedp.ScrollIntoView(`footer`, chromedp.WaitFunc(func(ctx context.Context, frame *cdp.Frame, ctxID runtime.ExecutionContextID, ids ...cdp.NodeID) ([]*cdp.Node, error) {

			// collecting data from web page
			err := collectSpeedResults(ctx, ch)
			if err != nil {
				return nil, fmt.Errorf("failed to collect results from web page: %w", err)
			}

			// routine for chromedp to finish method correctly
			nodes := make([]*cdp.Node, len(ids))
			frame.RLock()
			for i, id := range ids {
				nodes[i] = frame.Nodes[id]
				if nodes[i] == nil {
					frame.RUnlock()
					// not yet ready
					return nil, nil
				}
			}
			frame.RUnlock()
			return nodes, nil
		})),
	}

	err := chromedp.Run(ctx, cmds...)
	return err
}

func readUpdatesFromWebPage(ctx context.Context, succeededID, valueID, unitID string) (value, unit string, err error) {
	waitCtx, cancelFn := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancelFn()
	err = chromedp.WaitVisible(succeededID).Do(waitCtx)
	if err == nil {
		return "", "", errUpdatesFinished
	}

	chromedp.Text(valueID, &value, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
	chromedp.Text(unitID, &unit, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
	return value, unit, err
}

func collectSpeedResults(ctx context.Context, ch chan<- speed.Measures) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var (
		measures    speed.Measures
		value, unit string
		err         error
	)

	measures.Download = "*"
	ch <- measures

	// to skip first loop condition
	err = context.DeadlineExceeded
	for err != nil {
		<-ticker.C
		value, unit, err = readUpdatesFromWebPage(ctx, "#speed-value.succeeded", "#speed-value", "#speed-units")
		if err == errUpdatesFinished {
			err = nil
			break
		}
		measures.Download, measures.DownloadUnit = value, unit
		ch <- measures
	}

	measures.Upload = "*"
	ch <- measures

	err = chromedp.Click(`#show-more-details-link`).Do(ctx)
	if err != nil {
		return errors.New("uploading failed: no button found")
	}
	err = chromedp.WaitVisible(`#upload-value`).Do(ctx)
	if err != nil {
		return errors.New("uploading failed: no data found")
	}

	err = context.DeadlineExceeded
	for err != nil {
		<-ticker.C
		value, unit, err = readUpdatesFromWebPage(ctx, "#upload-value.succeeded", "#upload-value", "#upload-units")
		if err == errUpdatesFinished {
			err = nil
			break
		}
		measures.Upload, measures.UploadUnit = value, unit
		ch <- measures
	}

	return nil
}
