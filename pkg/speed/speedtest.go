package speed

import (
	"context"
	"fmt"
	"sync"
)

// Provider is an common interface for speed test providers
type Provider interface {
	Test(context.Context, chan<- Measures) error
}

// ProviderTest runs speed test on chosen provider and returns final speed measures
func ProviderTest(provider Provider) (measures Measures, err error) {
	ch := make(chan Measures, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = provider.Test(context.Background(), ch)
	}()

	for measures = range ch {
	}

	wg.Wait()

	return measures, err
}

// Test runs speed test on chosen provider and returns speed measures into the result channel whenever they are available
func Test(ctx context.Context, provider Provider, resultCh chan<- Measures) error {
	if provider == nil {
		return fmt.Errorf("nil speedtest provider")
	}
	if resultCh == nil {
		return fmt.Errorf("nil result channel")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return provider.Test(ctx, resultCh)
}
