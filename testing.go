package runner

import (
	"context"
	"encoding/json"
	"io"
)

// TestingT is a interface that describes the *testing.T methods needed by the
// Testing runner implementation.
type TestingT interface {
	Logf(format string, args ...interface{})
}

// Testing is a Runner that wraps another Runner, and logs all executed commands
// and their arguments to a *testing.T instance.
//
// Both Runner and T must be non-nil, or running commands will cause a panic.
type Testing struct {
	// Runner is the underlying Runner to run commands with. If not set, running
	// commands will cause a panic.
	Runner Runner

	// TestingT is the *testing.T instance used to log output. If not set,
	// running commands will cause a panic.
	TestingT TestingT

	// LogEnv indicates if calls to Env() should be logged.
	LogEnv bool
}

var _ Runner = &Testing{}

// Run executes the command with the underlying Runner, and logs command and
// arguments to TestingT.
func (r *Testing) Run(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	command string,
	args ...string,
) error {
	jsonArgs, _ := json.Marshal(args)
	r.TestingT.Logf(
		"runner.Run: command=%s args=%s",
		command, string(jsonArgs),
	)

	return r.Runner.Run(stdin, stdout, stderr, command, args...)
}

// RunContext executes the command with the underlying Runner, and logs command
// and arguments to TestingT.
func (r *Testing) RunContext(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	command string,
	args ...string,
) error {
	jsonArgs, _ := json.Marshal(args)
	r.TestingT.Logf(
		"runner.RunContext: command=%s args=%s",
		command, string(jsonArgs),
	)

	return r.Runner.RunContext(ctx, stdin, stdout, stderr, command, args...)
}

// Env sets the environment variables for the underlying Runner, and if LogEnv
// is true it logs the given environment variables to TestingT.
func (r *Testing) Env(vars ...string) {
	if r.LogEnv {
		jsonVars, _ := json.Marshal(vars)
		r.TestingT.Logf("runner.Env: vars=%s", string(jsonVars))
	}

	r.Runner.Env(vars...)
}
