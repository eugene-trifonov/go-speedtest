package speed

import (
	"fmt"
)

const (
	BitsPerSecond  = "bps"
	KbitsPerSecond = "Kbps"
	MbitsPerSecond = "Mbps"
	GbitsPerSecond = "Gbps"
)

var units = []string{BitsPerSecond, KbitsPerSecond, MbitsPerSecond, GbitsPerSecond}

// Format returns human readable values and units for internet speed
func Format(bytesPerSecond float64) (value, unit string) {
	// to make it bits per second
	numPerSecond := bytesPerSecond * 8
	i := 0
	for ; i < 3; i++ {
		if numPerSecond < 1024 {
			break
		}
		numPerSecond /= 1024
	}
	return fmt.Sprintf("%.2f", numPerSecond), units[i]
}
