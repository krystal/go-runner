package runner

import "io"

// Sudo is a Runner that wraps another Runner and runs commands via sudo.
//
// Password prompts are not supported, hence commands must be set to NOPASS via
// the sudoers file before they can be run.
type Sudo struct {
	// Runner is the underlying Runner to run commands with, after wrapping them
	// with sudo. Must be set, or Run will panic.
	Runner Runner

	// User value passed to sudo via -u flag.
	User string

	// Args is a string slice of extra arguments to pass to sudo.
	Args []string
}

var _ Runner = &Sudo{}

// Run executes the command via sudo by calling Run on the underlying Runner.
// Will panic if Runner field is nil on Sudo instance.
func (r *Sudo) Run(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	command string,
	args ...string,
) error {
	sudoArgs := []string{"-n"}
	if r.User != "" {
		sudoArgs = append(sudoArgs, "-u", r.User)
	}
	sudoArgs = append(sudoArgs, r.Args...)
	sudoArgs = append(sudoArgs, "--", command)
	sudoArgs = append(sudoArgs, args...)

	return r.Runner.Run(stdin, stdout, stderr, "sudo", sudoArgs...)
}

// Env sets the environment by calling Env on the underlying Runner. Will panic
// if Runner field is nil on Sudo instance.
func (r *Sudo) Env(vars ...string) {
	r.Runner.Env(vars...)
}
