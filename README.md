# go-speedtest
This library tests speed of an internet connection basing on third-party speed test providers using Netflix and Ookla.

Usage
```
import "github.com/go-speedtest/pkg/speed"
import "github.com/go-speedtest/pkg/netflix"

public main() {
	measures, err := speed.ProviderTest(netflix.Provider)
}
```

It allows to see the progress of speed test, to do so just pass channel:

Usage
```
import "github.com/go-speedtest/pkg/speed"
import "github.com/go-speedtest/pkg/netflix"

public main() {
	resultCh := make(chan speed.Measures)
	go func() {
		err := speed.Test(context.Background(), netflix.Provider, resultCh)
	}

	for measures := range resultCh {
		...
	}
}
```

The library can be run using command line
```
go run cmd/main.go ookla
speedtest.net: Download speed:   18.35  Mbps, Upload speed:   15.5  Mbps

go run cmd/main.go fast.com
fast.com: Download speed:   19.6  Mbps, Upload speed:   17.4  Mbps
```