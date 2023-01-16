package cmd

import (
	"flag"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"os"
)

type httpConfig struct {
	url    string
	verb   string
	output string
}

func HandleHttp(w io.Writer, args []string) error {
	var v string
	var o string
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&v, "verb", "GET", "HTTP method")
	fs.StringVar(&o, "output", "STDOUT", "Output Format")

	possibleVerb := []string{"GET", "POST", "HEAD"}
	// err 이면 0 이외에 종료코드 이겠찌...?
	if !slices.Contains(possibleVerb, v) {
		return InvalidHttpMethod
	}

	possibleOutput := []string{"STDOUT", "html"}
	if !slices.Contains(possibleOutput, o) {
		return UnsupportedOutputFormat
	}

	fs.Usage = func() {
		var usageString = `
http: A HTTP client.

http: <options> server`
		fmt.Fprintf(w, usageString)

		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	c := httpConfig{verb: v, output: o}
	c.url = fs.Arg(0)
	body, err := fetchRemoteResource(c.url)
	if err != nil {
		return err
	}

	if c.output == "STDOUT" {
		fmt.Fprintln(w, string(body))
	} else if c.output == "html" {
		f, err := os.OpenFile("output.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		f.Write(body)
	}

	return nil
}

func fetchRemoteResource(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
