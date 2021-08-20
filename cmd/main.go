package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-speedtest/pkg/netflix"
	"github.com/go-speedtest/pkg/ookla"
	"github.com/go-speedtest/pkg/speed"
)

var (
	providers = []string{ookla.Provider.String(), ookla.Name, netflix.Provider.String(), netflix.Name}
)

func main() {
	var provider speed.Provider = ookla.Provider
	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case ookla.Provider.String(), ookla.Name:
			provider = ookla.Provider
		case netflix.Provider.String(), netflix.Name:
			provider = netflix.Provider
		default:
			fmt.Printf("Please choose a provider to run the speed test:\n %s [%s]\n", os.Args[0], strings.Join(providers, "|"))
			return
		}
	}

	resultCh := make(chan speed.Measures, 1)
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	go func() {
		err := speed.Test(ctx, provider, resultCh)
		if err != nil {
			fmt.Println()
			log.Fatalf("failed to run speed test: %v", err)
		}
	}()

	running := true
	for running {
		select {
		case measures, ok := <-resultCh:
			if !ok {
				running = false
				break
			}
			fmt.Printf("%s: Download speed: %5s %5s, Upload speed: %5s %5s\r",
				provider, measures.Download, measures.DownloadUnit, measures.Upload, measures.UploadUnit)

		case <-ctx.Done():
			running = false

			// defer - to be sure it won't overwrite speed results
			defer fmt.Println("Sorry, we reached timeout")
		}
	}
	// to leave the results on terminal, otherwise it's overwritten because of \r
	fmt.Println()
}
