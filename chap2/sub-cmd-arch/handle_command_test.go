// Listing 2.6: chap2/sub-sub-cmd-arch/handle_command_test.go
package main

import (
	"bytes"
	"testing"
)

func TestHandleCommand(t *testing.T) {
	usageMessage := `Usage: mync [http|grpc] -h

http: A HTTP client.

http: <options> server

Options: 
  -verb string
    	HTTP method (default "GET")

grpc: A gRPC client.

grpc: <options> server

Options: 
  -body string
    	Body of request
  -method string
    	Method to call
`
	SampleUrl := "https://golang.org/pkg/net/http/"
	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		{
			args:   []string{},
			err:    errInvalidSubCommand,
			output: "Invalid sub-command specified\n" + usageMessage,
		},
		{
			args:   []string{"-h"},
			err:    nil,
			output: usageMessage,
		},
		{
			args:   []string{"foo"},
			err:    errInvalidSubCommand,
			output: "Invalid sub-command specified\n" + usageMessage,
		},
		{ // Exercise 2.1 (?) 왜 연습 문제 1.1의 해결 방법이 관련이 있는걸까?
			args:   []string{"http", SampleUrl},
			output: "Executing http command\n",
		},
		{
			args:   []string{"grpc", SampleUrl},
			output: "Executing grpc command\n",
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := handleCommand(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error %v, got %v", tc.err, err)
		}

		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Errorf("Expected output to be: %#v, Got: %#v", tc.output, gotOutput)
			}
		}
		byteBuf.Reset()
	}
}
