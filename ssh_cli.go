package runner

import (
	"context"
	"fmt"
	"io"
	"strconv"
)

var (
	ErrSSHCLI              = fmt.Errorf("%w: sshcli: ", Err)
	ErrSSHCLINoDestination = fmt.Errorf(
		"%w: destination must be set", ErrSSHCLI,
	)
)

// SSHCLI is a Runner that wraps another Runner, essentially prefixing given
// commands and arguments with "ssh", relevant SSH CLI arguments, and the given
// destination. It then passes this new "ssh" command to the underlying Runner.
//
// This is useful for running commands on remote hosts via SSH, without having
// to use the Go ssh package.
//
// Interactive commands are not supported, meaning SSH password prompts will not
// work, and the remote machine's hostkey should already be known and trusted by
// the ssh CLI client.
type SSHCLI struct {
	// Runner is the underlying Runner to run commands with, after wrapping them
	// with ssh. If not set, running commands will cause a panic.
	Runner Runner

	// Destination is the remote SSH destination to connect to, which may be
	// specified as either "[user@]hostname" or a URI of the form
	// "ssh://[user@]hostname[:port]".
	Destination string

	// Port is the remote SSH port (-p) flag to use. When 0, no -p flag will be
	// used.
	Port int

	// IdentityFile is the remote SSH identity file (-i) flag to use. When
	// empty, no -i flag will be used.
	IdentityFile string

	// Login is the remote SSH login (-l) flag to use. When empty, no -l flag
	// will be used.
	Login string

	// Args is a string slice of extra arguments to pass to ssh.
	Args []string

	env []string
}

var _ Runner = &SSHCLI{}

// Run executes the command remotely via ssh by calling Run on the underlying
// Runner.
//
// Will panic if Runner field is nil.
// Will return a error if Destination field is empty.
func (rsc *SSHCLI) Run(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	command string,
	args ...string,
) error {
	sshArgs, err := rsc.args(command, args)
	if err != nil {
		return err
	}

	return rsc.Runner.Run(stdin, stdout, stderr, "ssh", sshArgs...)
}

// RunContext executes the command remotely via ssh by calling RunContext on the
// underlying Runner.
//
// Will panic if Runner field is nil.
// Will return a error if Destination field is empty.
func (rsc *SSHCLI) RunContext(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	command string,
	args ...string,
) error {
	sshArgs, err := rsc.args(command, args)
	if err != nil {
		return err
	}

	return rsc.Runner.RunContext(ctx, stdin, stdout, stderr, "ssh", sshArgs...)
}

func (rsc *SSHCLI) args(command string, args []string) ([]string, error) {
	if rsc.Destination == "" {
		return nil, ErrSSHCLINoDestination
	}

	sshArgs := []string{}

	if rsc.Port != 0 {
		sshArgs = append(sshArgs, "-p", strconv.Itoa(rsc.Port))
	}
	if rsc.IdentityFile != "" {
		sshArgs = append(sshArgs, "-i", rsc.IdentityFile)
	}
	if rsc.Login != "" {
		sshArgs = append(sshArgs, "-l", rsc.Login)
	}
	if len(rsc.Args) > 0 {
		sshArgs = append(sshArgs, rsc.Args...)
	}
	sshArgs = append(sshArgs, rsc.Destination, "--")

	if len(rsc.env) > 0 {
		sshArgs = append(sshArgs, "env")
		sshArgs = append(sshArgs, rsc.env...)
	}
	sshArgs = append(sshArgs, command)
	sshArgs = append(sshArgs, args...)

	return sshArgs, nil
}

// Env sets the environment by calling Env on the underlying Runner. Will panic
// if Runner field is nil on SSH instance.
func (rsc *SSHCLI) Env(env ...string) {
	rsc.env = env
}
