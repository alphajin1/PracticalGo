package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
)

type config struct {
	numTimes int
	outPath  string
}

func getName(r io.Reader, w io.Writer) (string, error) {
	scanner := bufio.NewScanner(r)
	msg := "Your name please? Press the Enter key when done.\n"
	fmt.Fprintf(w, msg)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("You didn't enter your name")
	}
	return name, nil
}

func greetUser(c config, name string, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you %s\n", name)
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprintf(w, msg)
	}
}

func greetUserHTML(c config, name string, w io.Writer) {
	f, err := os.Create(c.outPath)
	if err != nil {
		errors.New("파일 생성 실패")
	}

	tmpl := `
				<h1>Nice to meet you {{.Name}}</h1>
			`
	t, err := template.New(c.outPath).Parse(tmpl)
	if err != nil {
		errors.New("template 생성 실패")
	}

	output := struct {
		Name string
	}{
		Name: name,
	}

	for i := 0; i < c.numTimes; i++ {
		t.Execute(f, output)
	}

}

func runCmd(r io.Reader, w io.Writer, c config) error {
	name, err := getName(r, w)
	if err != nil {
		return err
	}

	if c.outPath == "" {
		greetUser(c, name, w)
	} else {
		greetUserHTML(c, name, w)
	}

	return nil
}

func validateArgs(c config) error {
	if !(c.numTimes > 0) {
		return errors.New("Must specify a number greater than 0")
	}
	return nil
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("greeter", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.IntVar(&c.numTimes, "n", 0, "Number of times to greet")
	fs.StringVar(&c.outPath, "o", "", "생성할 html 파일 경로")
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}
	if fs.NArg() != 0 {
		return c, errors.New("Positional arguments specified")
	}
	return c, nil
}

func main() {
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	err = validateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	err = runCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
