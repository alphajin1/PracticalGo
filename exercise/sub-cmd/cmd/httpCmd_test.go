package cmd

import (
	"bytes"
	"testing"
)

func TestHandlePostCommand(t *testing.T) {

	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		{
			args: []string{"-verb", "GET", "https://golang.org/pkg/net/http/"},
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleHttp(byteBuf, tc.args)
		if err != nil {
			t.Fatalf("Error Exist!")
		}
	}
}
