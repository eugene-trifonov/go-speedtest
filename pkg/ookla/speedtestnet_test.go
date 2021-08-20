package ookla

import (
	"testing"

	"github.com/go-speedtest/pkg/speed"
)

func BenchmarkProviderTest(b *testing.B) {
	for n := 0; n < b.N; n++ {
		speed.ProviderTest(Provider)
	}
}
