package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestHandlePostCommand(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()

	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		{
			args: []string{"-verb", "POST", "-url", ts.URL, "-upload", "a.txt", "-form-data", "name=MyName", "-form-data", "version=1.0"},
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
