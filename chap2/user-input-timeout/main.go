package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var totalDuration time.Duration = 5

func getName(r io.Reader, w io.Writer) (string, error) {
	scanner := bufio.NewScanner(r)
	msg := "Your name please? Press the Enter key when done"
	fmt.Fprintln(w, msg)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("You entered an empty name")
	}
	return name, nil
}

func getNameContext(ctx context.Context) (string, error) {
	var err error
	name := "Default Name"
	c := make(chan error, 1)

	go func() {
		name, err = getName(os.Stdin, os.Stdout)
		// 에러값을 채널에 Write
		c <- err
	}()

	select {
	// 어떤 채널이 먼저 반응하는지에 따라 handling
	case <-ctx.Done():
		// ctx 가 종료되는 시점에 호출됨
		return name, ctx.Err()
	case err := <-c:
		// getName 함수가 반환할 때 쓰는 채널
		return name, err
	}
}

func main() {
	allowedDuration := totalDuration * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), allowedDuration)
	// ctx 해제를 위해 호출하는 것
	defer cancel()

	name, err := getNameContext(ctx)

	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, name)
}
