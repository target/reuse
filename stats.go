package main

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"log"
)

type Stats struct {
	DNSIntervals      stats.Float64Data
	TCPIntervals      stats.Float64Data
	SSLIntervals      stats.Float64Data
	ConnectIntevals   stats.Float64Data
	ServerIntervals   stats.Float64Data
	TransferIntervals stats.Float64Data
	TotalIntervals    stats.Float64Data
}

func (s *Stats) AddMetric(m *Metrics) {
	s.DNSIntervals = append(s.DNSIntervals, m.dnsTime())
	s.TCPIntervals = append(s.TCPIntervals, m.tcpTime())
	s.SSLIntervals = append(s.SSLIntervals, m.sslTime())
	s.ConnectIntevals = append(s.ConnectIntevals, m.connectTime())
	s.ServerIntervals = append(s.ServerIntervals, m.serverTime())
	s.TransferIntervals = append(s.TransferIntervals, m.transferTime())
	s.TotalIntervals = append(s.TotalIntervals, m.totalTime())
}

func (s *Stats) PrintStats() {
	fmt.Printf("--- Summary Statistics ---\n")
	fmt.Printf("         Min    Max    Mean   StdDev  Median  MAD\n")
	d := calculateStats(s.DNSIntervals)
	fmt.Printf("%-8s %-6.1f %-6.1f %-6.1f %-7.1f %-7.1f %-7.1f\n",
		"DNS", d[0], d[1], d[2], d[3], d[4], d[5])
	d = calculateStats(s.TCPIntervals)
	fmt.Printf("%-8s %-6.1f %-6.1f %-6.1f %-7.1f %-7.1f %-7.1f\n",
		"TCP", d[0], d[1], d[2], d[3], d[4], d[5])
	d = calculateStats(s.SSLIntervals)
	fmt.Printf("%-8s %-6.1f %-6.1f %-6.1f %-7.1f %-7.1f %-7.1f\n",
		"SSL", d[0], d[1], d[2], d[3], d[4], d[5])
	d = calculateStats(s.ConnectIntevals)
	fmt.Printf("%-8s %-6.1f %-6.1f %-6.1f %-7.1f %-7.1f %-7.1f\n",
		"TConn", d[0], d[1], d[2], d[3], d[4], d[5])
	d = calculateStats(s.ServerIntervals)
	fmt.Printf("%-8s %-6.1f %-6.1f %-6.1f %-7.1f %-7.1f %-7.1f\n",
		"Srv", d[0], d[1], d[2], d[3], d[4], d[5])
	d = calculateStats(s.TransferIntervals)
	fmt.Printf("%-8s %-6.1f %-6.1f %-6.1f %-7.1f %-7.1f %-7.1f\n",
		"Reply", d[0], d[1], d[2], d[3], d[4], d[5])
	d = calculateStats(s.TotalIntervals)
	fmt.Printf("%-8s %-6.1f %-6.1f %-6.1f %-7.1f %-7.1f %-7.1f\n",
		"Total", d[0], d[1], d[2], d[3], d[4], d[5])

}

func calculateStats(f stats.Float64Data) (results []float64) {
	min, err := f.Min()
	if err != nil {
		log.Fatalf("Error calculating Min: %v\n", err)
	}
	results = append(results, min)

	max, err := f.Max()
	if err != nil {
		log.Fatalf("Error calculating Max: %v\n", err)
	}
	results = append(results, max)

	mean, err := f.Mean()
	if err != nil {
		log.Fatalf("Error calculating Mean: %v\n", err)
	}
	results = append(results, mean)

	stddev, err := f.StandardDeviation()
	if err != nil {
		log.Fatalf("Error calculating StdDev: %v\n", err)
	}
	results = append(results, stddev)

	median, err := f.Median()
	if err != nil {
		log.Fatalf("Error calculating Median: %v\n", err)
	}
	results = append(results, median)

	mad, err := f.MedianAbsoluteDeviation()
	if err != nil {
		log.Fatalf("Error calculating MAD: %v\n", err)
	}
	results = append(results, mad)

	return results
}
