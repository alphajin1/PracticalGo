package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestHandlePostCommand(t *testing.T) {

	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		{
			args: []string{"-verb", "GET", "-url", "https://golang.org/pkg/net/http/"},
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleHttp(byteBuf, tc.args)
		if err != nil {
			t.Fatalf("Error Exist!, %s", err)
		}

		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if !strings.Contains(gotOutput, tc.output) {
				t.Errorf("Expected output to be: %#v, Got: %#v", tc.output, gotOutput)
			}
		}
	}
}
