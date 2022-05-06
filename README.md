<h1 align="center">
  go-runner
</h1>

<p align="center">
  <strong>
    Go package that exposes a `Runner` interface for executing commands locally
    via exec.Command.
  </strong>
</p>

<p align="center">
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
</p>

## Import

```go
import "github.com/krystal/go-runner"
```

## Usage

```go
var stdout, stderr bytes.Buffer
stdin := bytes.NewBufferString("Hello world!")

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

Output:

```
Oh noes! :(
Hello world!
```
