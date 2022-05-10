<h1 align="center">
  go-runner
</h1>

<p align="center">
  <strong>
    Go package exposing a simple interface for executing commands, enabling easy
    mocking and wrapping of executed commands.
  </strong>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/krystal/go-runner">
    <img src="https://img.shields.io/badge/%E2%80%8B-reference-387b97.svg?logo=go&logoColor=white"
  alt="Go Reference">
  </a>
  <a href="https://github.com/krystal/go-runner/releases">
    <img src="https://img.shields.io/github/v/tag/krystal/go-runner?label=release" alt="GitHub tag (latest SemVer)">
  </a>
  <a href="https://github.com/krystal/go-runner/actions">
    <img src="https://img.shields.io/github/workflow/status/krystal/go-runner/CI.svg?logo=github" alt="Actions Status">
  </a>
  <a href="https://github.com/krystal/go-runner/commits/main">
    <img src="https://img.shields.io/github/last-commit/krystal/go-runner.svg?style=flat&logo=github&logoColor=white"
alt="GitHub last commit">
  </a>
  <a href="https://github.com/krystal/go-runner/issues">
    <img src="https://img.shields.io/github/issues-raw/krystal/go-runner.svg?style=flat&logo=github&logoColor=white"
alt="GitHub issues">
  </a>
  <a href="https://github.com/krystal/go-runner/pulls">
    <img src="https://img.shields.io/github/issues-pr-raw/krystal/go-runner.svg?style=flat&logo=github&logoColor=white" alt="GitHub pull requests">
  </a>
  <a href="https://github.com/krystal/go-runner/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/krystal/go-runner.svg?style=flat" alt="License Status">
  </a>
</p>

The Runner interface is basic and minimal, but it is sufficient for most use
cases. This makes it easy to mock Runner for testing purposes.

It's also easy to create wrapper runners which modify commands before executing
them. The `Sudo` struct is a simple example of this.

## Import

```go
import "github.com/krystal/go-runner"
```

## Interface

```go
type Runner interface {
	Run(
		stdin io.Reader,
		stdout, stderr io.Writer,
		command string,
		args ...string,
	) error
	RunContext(
		ctx context.Context,
		stdin io.Reader,
		stdout, stderr io.Writer,
		command string,
		args ...string,
	) error
	Env(env ...string)
}
```

## Usage

Basic:

```go
var stdout bytes.Buffer

r := runner.New()
_ = r.Run(nil, &stdout, nil, "echo", "Hello world!")

fmt.Print(stdout.String())
```

```
Hello world!
```

Environment:

```go
var stdout bytes.Buffer

r := runner.New()
r.Env("USER=johndoe", "HOME=/home/johnny")
_ = r.Run(nil, &stdout, nil, "sh", "-c", `echo "Hi, ${USER} (${HOME})"`)

fmt.Print(stdout.String())
```

```
Hi, johndoe (/home/johnny)
```

Stdin, Stdout, and Stderr:

```go
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
```

```
Oh noes! :(
Hello world!
```

Failure:

```go
var stdout, stderr bytes.Buffer

r := runner.New()
err := r.Run(
	nil, &stdout, &stderr,
	"sh", "-c", "echo 'Hello world!'; echo 'Oh noes! :(' >&2; exit 3",
)
if err != nil {
	fmt.Printf("%s: %s", err.Error(), stderr.String())
}
```

```
exit status 3: Oh noes! :(
```

Context:

```go
var stdout bytes.Buffer

ctx, cancel := context.WithTimeout(
	context.Background(), 1*time.Second,
)
defer cancel()

r := runner.New()
err := r.RunContext(
	ctx, nil, &stdout, nil,
	"sh", "-c", "sleep 0.5 && echo 'Hello world!'",
)
if err != nil {
	fmt.Println(err)
}

fmt.Print(stdout.String())
```

```
Hello world!
```

Sudo (requires `NOPASS` in sudoers file):

```go
var stdout bytes.Buffer
r := runner.New()

sudo := &runner.Sudo{Runner: r}
_ = sudo.Run(nil, &stdout, nil, "whoami")

sudo.User = "web"
_ = sudo.Run(nil, &stdout, nil, "whoami")

fmt.Print(stdout.String())
```

```
root
web
```

## Documentation

Please see the
[Go Reference](https://pkg.go.dev/github.com/krystal/go-runner#section-documentation)
for documentation and examples.

## License

[MIT](https://github.com/krystal/go-runner/blob/main/LICENSE)
