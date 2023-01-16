package cmd

import (
	"flag"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
)

type HttpStatus struct {
	code       int
	message    string
	customCode int
}

type httpConfig struct {
	url    string
	verb   string
	body   string
	output string
}

func isPossibleConfig(fs *flag.FlagSet, c *httpConfig) (bool, error) {
	possibleVerb := []string{"GET", "POST", "HEAD"}
	if !slices.Contains(possibleVerb, c.verb) {
		return false, UnSupportedHTTPMethod
	}

	possibleOutput := []string{"STDOUT", "html"}
	if !slices.Contains(possibleOutput, c.output) {
		return false, UnSupportedOutputFormat
	}

	if fs.NArg() != 1 {
		return false, ErrNoServerSpecified
	}

	return true, nil
}

func HandleHttp(w io.Writer, args []string) (HttpStatus, error) {
	c := httpConfig{}
	r := HttpStatus{}
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.verb, "verb", "GET", "HTTP method")
	fs.StringVar(&c.body, "body", "{}", "Body")
	fs.StringVar(&c.output, "output", "STDOUT", "Output Format")

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
		return r, err
	}

	_, err = isPossibleConfig(fs, &c)
	if err != nil {
		return r, err
	}

	c.url = fs.Arg(0)

	if c.verb == "GET" {
		body, err := fetchRemoteResource(c.url)
		if err != nil {
			return r, err
		}

		err = flushOutput(w, body, c.output)
		if err != nil {
			return r, err
		}
	}

	if c.verb == "POST" {

	}

	return r, nil
}
