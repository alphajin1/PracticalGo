package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	tests := []struct {
		arg           string
		name          string
		containString string
		exitString    string
	}{
		{
			arg:           "2",
			name:          "Honey",
			containString: "Your name please? Press the Enter key when done.\n" + strings.Repeat("Nice to meet you Honey\n", 2),
		},
		{
			arg:           "0",
			containString: "Must specify a number greater than 0\n",
			exitString:    fmt.Sprintf("exit status %d", 1),
		},
		{
			containString: "Usage: ./application <integer> [-h|-help]",
			exitString:    fmt.Sprintf("exit status %d", 1),
		},
	}

	for _, tc := range tests {
		cmd := exec.Command("./application", tc.arg)
		cmd.Stdin = strings.NewReader(tc.name)

		out, err := cmd.CombinedOutput()
		sout := string(out)
		if !strings.Contains(sout, tc.containString) {
			log.Fatalf("contains string is not matched. %s != %s", sout, tc.containString)
		}
		if err != nil && err.Error() != tc.exitString {
			log.Fatalf("exit status is not matched. %s != %s", err.Error(), tc.exitString)
		}
	}
}
