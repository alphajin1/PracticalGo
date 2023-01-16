package cmd

import (
	"flag"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"strings"
)

// Custom Arguments
type formDataList []string

func (i *formDataList) String() string {
	return "MyStringRepresentation"
}

func (i *formDataList) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type httpConfig struct {
	url      string
	verb     string // http Types
	output   string // GET
	upload   string // POST, ex) /path/to/file.pdf
	formData formDataList
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

	return true, nil
}

func HandleHttp(w io.Writer, args []string) error {
	c := httpConfig{}
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.url, "url", "", "Request URL")
	fs.StringVar(&c.verb, "verb", "GET", "HTTP method")
	fs.StringVar(&c.output, "output", "STDOUT", "Output Format")
	fs.StringVar(&c.upload, "upload", "", "Upload FileName")
	fs.Var(&c.formData, "form-data", "Request Form Data")
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

	_, err = isPossibleConfig(fs, &c)
	if err != nil {
		return err
	}

	if c.verb == "GET" {
		body, err := fetchRemoteResource(c.url)
		if err != nil {
			return err
		}

		err = flushOutput(w, body, c.output)
		if err != nil {
			return err
		}
	}

	if c.verb == "POST" {
		m := make(map[string]string)
		for _, v := range c.formData {
			parts := strings.Split(v, "=")
			m[parts[0]] = parts[1]
		}

		p := pkgData{
			Name:     m["name"],
			Version:  m["version"],
			Filename: c.upload,
			Bytes:    strings.NewReader("data"),
		}

		response, err := registerPackageData(c.url, p)
		if err != nil {
			return err
		}

		err = flushOutput(w, []byte(response.ID), c.output)
		if err != nil {
			return err
		}
	}

	return nil
}
