package main

import (
	"fmt"
	"strings"
	"time"
)

type Metrics struct {
	TCPReused                bool
	SSLReused                bool
	DNSCoalesced             bool
	SSL                      bool
	HTTPStatus               int
	LocalAddr                string
	RemoteAddr               string
	StartTime                time.Time
	GetConnTime              time.Time
	GotConnTime              time.Time
	GotFirstResponseByteTime time.Time
	DNSStartTime             time.Time
	DNSDoneTime              time.Time
	ConnectStartTime         time.Time
	ConnectDoneTime          time.Time
	WroteRequestTime         time.Time
	EndTime                  time.Time
}

func PrintHeader() {
	fmt.Printf("lport Remote Address        Rsp DTS DNS   TCP   SSL   TConn Srv   Reply Total\n")
}

func (m Metrics) Print() {
	fmt.Printf("%-5s %-21s %-3d %d%d%d %-5.1f %-5.1f %-5.1f %-5.1f %-5.1f %-5.1f %-5.1f\n",
		m.localPort(),
		m.RemoteAddr,
		m.HTTPStatus,
		boolToInt(m.DNSCoalesced),
		boolToInt(m.TCPReused),
		boolToInt(m.SSLReused),
		m.dnsTime(),
		m.tcpTime(),
		m.sslTime(),
		m.connectTime(),
		m.serverTime(),
		m.transferTime(),
		m.totalTime())
}

func (m Metrics) localPort() string {
	s := strings.Split(m.LocalAddr, ":")
	return s[1]
}

func (m Metrics) dnsTime() float64 {
	return float64(m.DNSDoneTime.Sub(m.DNSStartTime).Nanoseconds()) / float64(1e6)
}

func (m Metrics) tcpTime() float64 {
	if m.ConnectDoneTime.IsZero() {
		return float64(0)
	}
	if m.DNSDoneTime.IsZero() {
		return float64(m.ConnectDoneTime.Sub(m.ConnectStartTime).Nanoseconds()) / float64(1e6)
	}
	return float64(m.ConnectDoneTime.Sub(m.DNSDoneTime).Nanoseconds()) / float64(1e6)
}

func (m Metrics) sslTime() float64 {
	if m.SSL && !m.ConnectDoneTime.IsZero() {
		return float64(m.WroteRequestTime.Sub(m.ConnectDoneTime).Nanoseconds()) / float64(1e6)
	} else {
		return float64(0)
	}
}

func (m Metrics) connectTime() float64 {
	return float64(m.GotConnTime.Sub(m.GetConnTime).Nanoseconds()) / float64(1e6)
}

func (m Metrics) serverTime() float64 {
	return float64(m.GotFirstResponseByteTime.Sub(m.WroteRequestTime).Nanoseconds()) / float64(1e6)
}
func (m Metrics) transferTime() float64 {
	return float64(m.EndTime.Sub(m.GotFirstResponseByteTime).Nanoseconds()) / float64(1e6)
}

func (m Metrics) totalTime() float64 {
	return float64(m.EndTime.Sub(m.StartTime).Nanoseconds()) / float64(1e6)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
