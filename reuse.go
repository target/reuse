package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strings"
	"time"
)

var version = "unversioned"

func main() {
	var stats Stats
	var opts Options

	_, err := flags.Parse(&opts)

	if opts.PrintVersion {
		fmt.Printf("reuse version %s.\n", version)
		os.Exit(0)
	}
	if err != nil {
		os.Exit(1)
	}

	ctx := context.Background()
	client := getClient(opts)
	PrintHeader()
	for i := 0; i < opts.Repetions; i++ {
		var metrics Metrics
		ctx := setTrace(ctx, &metrics)
		req, err := buildRequest(opts)
		if err != nil {
			log.Fatalf("Couldn't build request to %s:%v\n", opts.Args.Url, err)
		}
		req = req.WithContext(ctx)

		metrics.StartTime = time.Now()
		resp, err := client.Do(req)

		if err != nil {
			log.Fatalf("Couldn't GET %s:%v\n", opts.Args.Url, err)
		}

		if resp.TLS != nil {
			metrics.SSL = true
			metrics.SSLReused = resp.TLS.DidResume
		}

		metrics.HTTPStatus = resp.StatusCode

		_, err = ioutil.ReadAll(resp.Body)
		metrics.EndTime = time.Now()
		resp.Body.Close()

		if err != nil {
			log.Fatalf("reuse: error reading response body from %s:%v", opts.Args.Url, err)
		}

		metrics.Print()
		stats.AddMetric(&metrics)

		if i < (opts.Repetions - 1) {
			time.Sleep(opts.Wait)
		}
	}
	stats.PrintStats()
}

func getTransport(o Options) *http.Transport {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: o.Insecure,
			ClientSessionCache: tls.NewLRUClientSessionCache(25),
		},
	}
	if len(o.Proxy) > 0 {
		proxyUrl, _ := url.Parse(o.Proxy)
		tr.Proxy = http.ProxyURL(proxyUrl)
	}
	return tr
}

func getClient(o Options) *http.Client {
	t := getTransport(o)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= o.MaxRedirs {
				return http.ErrUseLastResponse
			}
			return nil
		},
		Transport: t}
	return client
}

func buildRequest(o Options) (req *http.Request, err error) {
	var reader io.Reader
	if o.Data != "" {
		reader = strings.NewReader(o.Data)
	}
	req, err = http.NewRequest(o.Method, o.Args.Url, reader)
	for _, v := range o.Headers {
		req.Header.Set(v.Name, v.Value)
	}

	return req, err
}

func setTrace(ctx context.Context, m *Metrics) context.Context {
	c := httptrace.WithClientTrace(ctx,
		&httptrace.ClientTrace{
			GetConn: func(hostPort string) {
				m.GetConnTime = time.Now()
			},
			GotConn: func(info httptrace.GotConnInfo) {
				m.GotConnTime = time.Now()
				m.TCPReused = info.Reused
				m.RemoteAddr = info.Conn.RemoteAddr().String()
				m.LocalAddr = info.Conn.LocalAddr().String()
			},
			GotFirstResponseByte: func() {
				m.GotFirstResponseByteTime = time.Now()
			},
			DNSStart: func(info httptrace.DNSStartInfo) {
				m.DNSStartTime = time.Now()
			},
			DNSDone: func(info httptrace.DNSDoneInfo) {
				m.DNSDoneTime = time.Now()
				m.DNSCoalesced = info.Coalesced
			},
			ConnectStart: func(network, addr string) {
				m.ConnectStartTime = time.Now()
			},
			ConnectDone: func(netowrk, addr string, err error) {
				m.ConnectDoneTime = time.Now()
			},
			WroteRequest: func(httptrace.WroteRequestInfo) {
				m.WroteRequestTime = time.Now()
			},
		})
	return c
}
