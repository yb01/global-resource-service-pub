package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GetLatencyReport(t *testing.T) {
	testCases := []struct {
		name  string
		start int
		gap   int
		count int
		p50   int
		p90   int
		p99   int
	}{
		{
			name:  "Test 1, start 1",
			start: 1,
			gap:   1,
			count: 1,
			p50:   1,
			p90:   1,
			p99:   1,
		},
		{
			name:  "Test 1-100, start 1",
			start: 1,
			gap:   1,
			count: 100,
			p50:   50,
			p90:   90,
			p99:   99,
		},
		{
			name:  "Test 1-99, start 1",
			start: 1,
			gap:   1,
			count: 99,
			p50:   50,
			p90:   90,
			p99:   99,
		},
		{
			name:  "Test 1-1000, start 10",
			start: 10,
			gap:   1,
			count: 1000,
			p50:   509,
			p90:   909,
			p99:   999,
		},
		{
			name:  "Test 1-1M, start 1",
			start: 1,
			gap:   1,
			count: 1000000,
			p50:   500000,
			p90:   900000,
			p99:   990000,
		},
		{
			name:  "Test 1-10M, start 1",
			start: 1,
			gap:   1,
			count: 10000000,
			p50:   5000000,
			p90:   9000000,
			p99:   9900000,
		},
		{
			name:  "Test 10M-1, start 10M",
			start: 10000000,
			gap:   -1,
			count: 10000000,
			p50:   5000000,
			p90:   9000000,
			p99:   9900000,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			start := tt.start

			latency := NewLatencyMetrics(0)
			for i := 0; i < tt.count; i++ {
				latency.AddLatencyMetrics(time.Duration(start) * time.Millisecond)
				start += tt.gap
			}

			report := latency.GetSummary()
			t.Logf("Test [%s], p50 %v, p90 %v, p99 %v, count %v", tt.name, report.P50, report.P90, report.P99, report.TotalCount)
			assert.Equal(t, time.Duration(tt.p50)*time.Millisecond, report.P50)
			assert.Equal(t, time.Duration(tt.p90)*time.Millisecond, report.P90)
			assert.Equal(t, time.Duration(tt.p99)*time.Millisecond, report.P99)
			assert.Equal(t, tt.count, report.TotalCount)
		})
	}
}
