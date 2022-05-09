package runner_test

import (
	"bytes"
	"fmt"

	"github.com/krystal/go-runner"
)

func ExampleRunner_basic() {
	var stdout bytes.Buffer

	r := runner.New()
	_ = r.Run(nil, &stdout, nil, "echo", "Hello world!")

	fmt.Print(stdout.String())
	// Output:
	// Hello world!
}

func ExampleRunner_environment() {
	var stdout bytes.Buffer

	r := runner.New()
	r.Env("USER=johndoe", "EMAIL=john@doe.io")
	_ = r.Run(nil, &stdout, nil, "sh", "-c", `echo "Hi, ${USER} (${EMAIL})"`)

	fmt.Print(stdout.String())
	// Output:
	// Hi, johndoe (john@doe.io)
}

func ExampleRunner_stdin() {
	var stdout bytes.Buffer

	r := runner.New()
	_ = r.Run(bytes.NewBufferString("Hello world!"), &stdout, nil, "cat")

	fmt.Print(stdout.String())
	// Output:
	// Hello world!
}

func ExampleRunner_stderr() {
	var stderr bytes.Buffer

	r := runner.New()
	err := r.Run(nil, nil, &stderr, "sh", "-c", "echo 'Oh noes! :(' >&2")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(stderr.String())
	// Output:
	// Oh noes! :(
}

func ExampleRunner_stdoutAndStderr() {
	var stdout, stderr bytes.Buffer

	r := runner.New()
	err := r.Run(
		nil, &stdout, &stderr,
		"sh", "-c", "echo 'Hello world!'; echo 'Oh noes! :(' >&2",
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(stderr.String())
	fmt.Print(stdout.String())
	// Output:
	// Oh noes! :(
	// Hello world!
}

func ExampleRunner_stdinStdoutAndStderr() {
	stdin := bytes.NewBufferString("Hello world!")
	var stdout, stderr bytes.Buffer

	r := runner.New()
	err := r.Run(
		stdin, &stdout, &stderr,
		"sh", "-c", "cat; echo 'Oh noes! :(' >&2",
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(stderr.String())
	fmt.Print(stdout.String())
	// Output:
	// Oh noes! :(
	// Hello world!
}

func ExampleRunner_combined() {
	var out bytes.Buffer

	r := runner.New()
	err := r.Run(
		nil, &out, &out,
		"sh", "-c", "echo 'Hello world!'; echo 'Oh noes! :(' >&2",
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(out.String())
	// Output:
	// Hello world!
	// Oh noes! :(
}

func ExampleRunner_failure() {
	var stdout, stderr bytes.Buffer

	r := runner.New()
	err := r.Run(
		nil, &stdout, &stderr,
		"sh", "-c", "echo 'Hello world!'; echo 'Oh noes! :(' >&2; exit 3",
	)
	if err != nil {
		fmt.Printf("%s: %s", err.Error(), stderr.String())
	}

	// Output:
	// exit status 3: Oh noes! :(
}
