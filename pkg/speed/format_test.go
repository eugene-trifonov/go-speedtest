package speed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name           string
		bytesPerSecond float64
		expectedValue  string
		expectedUnit   string
	}{
		{
			name:           "zero",
			bytesPerSecond: 0,
			expectedValue:  "0.00",
			expectedUnit:   BitsPerSecond,
		},
		{
			name:           "bytes to bits",
			bytesPerSecond: 1,
			expectedValue:  "8.00",
			expectedUnit:   BitsPerSecond,
		},
		{
			name:           "1Kbps",
			bytesPerSecond: 128,
			expectedValue:  "1.00",
			expectedUnit:   KbitsPerSecond,
		},
		{
			name:           "100Kbps",
			bytesPerSecond: 12800,
			expectedValue:  "100.00",
			expectedUnit:   KbitsPerSecond,
		},
		{
			name:           "1Mbps",
			bytesPerSecond: 128 * 1024,
			expectedValue:  "1.00",
			expectedUnit:   MbitsPerSecond,
		},
		{
			name:           "1.5Mbps",
			bytesPerSecond: 192 * 1024,
			expectedValue:  "1.50",
			expectedUnit:   MbitsPerSecond,
		},
		{
			name:           "1Gbps",
			bytesPerSecond: 2 * 1024 * 1024 * 1024 / 8,
			expectedValue:  "2.00",
			expectedUnit:   GbitsPerSecond,
		},
		{
			name:           "Gbps - max unit",
			bytesPerSecond: 1024 * 1024 * 1024 * 1024 / 8,
			expectedValue:  "1024.00",
			expectedUnit:   GbitsPerSecond,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, unit := Format(test.bytesPerSecond)
			assert.Equal(t, test.expectedValue, value)
			assert.Equal(t, test.expectedUnit, unit)
		})
	}
}
