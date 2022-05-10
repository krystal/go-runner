// Package runner exposes a simple interface for executing commands, enabling
// easy mocking and wrapping of executed commands.
//
// The Runner interface is basic and minimal, but it is sufficient for most use
// cases. This makes it easy to mock Runner for testing purposes.
//
// It's also easy to create wrapper runners which modify commands before
// executing them. The Sudo struct is a simple example of this.
package runner

import (
	"context"
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
		stdout, stderr io.Writer,
		command string,
		args ...string,
	) error

	// RunContext is like Run but includes a context.
	//
	// The provided context is used to kill the command process if the context
	// becomes done before the command completes on its own.
	RunContext(
		ctx context.Context,
		stdin io.Reader,
		stdout, stderr io.Writer,
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
	cmd := exec.Command(command, args...)

	return r.run(cmd, stdin, stdout, stderr)
}

// RunContext executes the given command locally on the host machine, using the
// provided context to kill the process if the context becomes done before the
// command completes on its own.
func (r *Local) RunContext(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	command string,
	args ...string,
) error {
	cmd := exec.CommandContext(ctx, command, args...)

	return r.run(cmd, stdin, stdout, stderr)
}

func (r *Local) run(
	cmd *exec.Cmd,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

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
