package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func fetchRemoteResource(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func flushOutput(w io.Writer, body []byte, output string) error {
	if output == "STDOUT" {
		fmt.Fprintln(w, string(body))
	} else if output == "html" {
		fName := fmt.Sprintf("output.%s", output)
		f, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		f.Write(body)
	}

	return nil
}
