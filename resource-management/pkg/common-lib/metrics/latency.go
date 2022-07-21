package metrics

import (
	"math"
	"sort"
	"time"
)

type LatencyReport struct {
	TotalCount int
	P50        time.Duration
	P90        time.Duration
	P99        time.Duration
}

type LatencyMetrics struct {
	value     int
	latencies []time.Duration
}

func NewLatencyMetrics(value int) *LatencyMetrics {
	return &LatencyMetrics{
		value:     value,
		latencies: make([]time.Duration, 0),
	}
}

func (m *LatencyMetrics) AddLatencyMetrics(newLatency time.Duration) {
	m.latencies = append(m.latencies, newLatency)
}

func (m *LatencyMetrics) Len() int {
	return len(m.latencies)
}

func (m *LatencyMetrics) Less(i, j int) bool {
	return m.latencies[i] < m.latencies[j]
}

func (m *LatencyMetrics) Swap(i, j int) {
	m.latencies[i], m.latencies[j] = m.latencies[j], m.latencies[i]
}

func (m *LatencyMetrics) GetSummary() *LatencyReport {
	// sort
	sort.Sort(m)
	count := len(m.latencies)
	if count == 0 {
		return &LatencyReport{
			TotalCount: count,
			P50:        0,
			P90:        0,
			P99:        0,
		}
	}
	return &LatencyReport{
		TotalCount: count,
		P50:        m.latencies[int(math.Ceil(float64(count*50)/100))-1],
		P90:        m.latencies[int(math.Ceil(float64(count*90)/100))-1],
		P99:        m.latencies[int(math.Ceil(float64(count*99)/100))-1],
	}
}
