package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func createContextWithTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	return ctx, cancel
}

func setupSignalHandler(w io.Writer, cancelFunc context.CancelFunc) {
	// 어느 한 시그널을 수신할 경우 시그널 핸들링 코드는 1단계(WithTimeout) 에서 반환된 취소 함수(cancel)을 호출
	c := make(chan os.Signal, 1)
	// SIGINT, SIGTERM에 대한 수신채널을 설정
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-c
		fmt.Fprintf(w, "Got signal: %v\n", s)
		cancelFunc()
	}()
}

func executeCommand(ctx context.Context, command string, arg string) error {
	return exec.CommandContext(ctx, command, arg).Run()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stdout, "Usage: %s <command> <argument>\n", os.Args[0])
		os.Exit(1)
	}
	command := os.Args[1]
	arg := os.Args[2]

	cmdTimeout := 30 * time.Second
	ctx, cancel := createContextWithTimeout(cmdTimeout)
	defer cancel()
	setupSignalHandler(os.Stdout, cancel)

	err := executeCommand(ctx, command, arg)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
