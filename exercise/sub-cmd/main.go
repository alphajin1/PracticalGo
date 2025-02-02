package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alphajin1/PracticalGo/exercise/sub-cmd/cmd"
)

var errInvalidSubCommand = errors.New("invalid sub-command specified")

func printUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: mync [http|grpc] -h\n")
	cmd.HandleHttp(w, []string{"-h"})
	cmd.HandleGrpc(w, []string{"-h"})
}

func handleCommand(w io.Writer, args []string) error {
	var err error

	if len(args) < 1 {
		// 인수로 아무 값도 전달되지 않은 경우
		err = errInvalidSubCommand
	} else {
		switch args[0] {
		case "http":
			err = cmd.HandleHttp(w, args[1:])
		case "grpc":
			err = cmd.HandleGrpc(w, args[1:])
		case "-h":
			printUsage(w)
		case "-help":
			printUsage(w)
		default:
			err = errInvalidSubCommand
		}
	}
	if errors.Is(err, cmd.ErrNoServerSpecified) || errors.Is(err, errInvalidSubCommand) {
		fmt.Fprintln(w, err)
		printUsage(w)
	}
	return err
}

func main() {
	err := handleCommand(os.Stdout, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	// ./application http https://golang.org/pkg/net/http/
	// [RESULT] Executing http command

	// ./application grpc https://golang.org/pkg/net/http/
	// [RESULT] Executing grpc command

	// ./application http -output html https://golang.org/pkg/net/http/
	// [RESULT] output.html file
}
