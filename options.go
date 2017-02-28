package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Header struct {
	Name  string
	Value string
}

func (h *Header) UnmarshalFlag(value string) error {
	parts := strings.Split(value, ":")

	if len(parts) != 2 {
		return errors.New("expected 2 strings separated by a :")
	}
	h.Name = parts[0]
	h.Value = parts[1]
	return nil
}

func (h Header) MarshalFlag() (string, error) {
	return fmt.Sprintf("%s:%s", h.Name, h.Value), nil
}

type Options struct {
	Repetions    int           `short:"r" long:"repetitions" description:"Number of times to repeat connecting" default:"2"`
	MaxRedirs    int           `long:"max-redirs" description:"Maximum number of redirects to follow" default:"10"`
	Insecure     bool          `short:"k" long:"insecure" description:"Skip SSL Verification"`
	PrintVersion bool          `short:"V" long:"version" description:"Print Version number and exit"`
	Data         string        `short:"d" long:"data" description:"Data to send"`
	Method       string        `short:"X" long:"request" description:"HTTP method" default:"GET"`
	Wait         time.Duration `short:"w" long:"wait" description:"Time to wait between connections" default:"5s"`
	Headers      []Header      `short:"H" long:"header" description:"Additional Header"`
	Args         struct {
		Url string `positional-arg-name:"URL" description:"URL to connect to"`
	} `positional-args:"yes" required:"yes"`
}
