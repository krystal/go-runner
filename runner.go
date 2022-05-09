// Package runner exposes a simple interface for executing commands, enabling
// easy mocking and wrapping of executed commands.
//
// It enables easy mocking of runners for testing purposes, and also for
// wrapping a runner to modify the commands being executed, like passing all
// commands through sudo for example.
package runner

import (
	"io"
	"os/exec"
)

//go:generate mockgen -source=$GOFILE -destination=mock/${GOFILE}

// Runner is the interface that Manager uses internally to run commands. This
// makes it easy to replace the underlying command runner with a mock for
// testing, or a different runner that executes givens commands in a different
// manner.
type Runner interface {
	// Run executes the given command with any provided arguments. Stdin,
	// Stdout, and Stderr can be provided/captured if the io.Reader/Writer is
	// not nil.
	Run(
		stdin io.Reader,
		stdout io.Writer,
		stderr io.Writer,
		command string,
		args ...string,
	) error

	// Env specifies the environment variables which will be available to all
	// commands invoked by the runner. Each entry is of the form "key=value".
	// Entries with duplicate keys will cause all but the last to be ignored.
	//
	// Multiple calls to Env will overwrite any previous calls to Env.
	//
	// If no env is set, no environment variables will be set for executed
	// commands.
	//
	// To set the environment to match that of the Go runtime, call Env with
	// os.Environ().
	Env(env ...string)
}

// Local is a Runner implementation that executes commands locally on the
// host machine.
type Local struct {
	env []string
}

var _ Runner = &Local{}

// New returns a Local instance which meets the Runner interface, and executes
// commands locally on the host machine.
func New() Runner {
	return &Local{}
}

// Run executes the given command locally on the host machine.
func (r *Local) Run(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	command string,
	args ...string,
) error {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = r.env
	if stdin != nil {
		cmd.Stdin = stdin
	}

	return cmd.Run()
}

// Env sets the environment which will apply to all commands invoked by the
// runner. Each entry is of the form "key=value".
func (r *Local) Env(env ...string) {
	r.env = env
}
