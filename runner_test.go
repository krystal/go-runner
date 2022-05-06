package runner

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r := New()

	assert.NotNil(t, r)
	assert.IsType(t, (*Local)(nil), r)
	assert.Implements(t, (*Runner)(nil), r)
}

func TestLocal_Run(t *testing.T) {
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "helloworld")
	require.NoError(t, err)
	_, err = f.WriteString("hello world :)")
	require.NoError(t, err)
	f.Close()
	helloFile := f.Name()

	tests := []struct {
		name          string
		stdin         []byte
		command       string
		args          []string
		discardStdout bool
		discardStderr bool
		wantStdout    []byte
		wantStderr    []byte
		wantErr       string
	}{
		{
			name:       "echo",
			command:    "echo",
			args:       []string{"hello", "world"},
			wantStdout: []byte("hello world\n"),
		},
		{
			name:       "cat file",
			command:    "cat",
			args:       []string{helloFile},
			wantStdout: []byte("hello world :)"),
		},
		{
			name:       "cat stdin",
			stdin:      []byte("this is some stdin text"),
			command:    "cat",
			wantStdout: []byte("this is some stdin text"),
		},
		{
			name: "cat multi-line stdin",
			stdin: []byte(`this is some stdin text
and some more text

and some more again :P
`),
			command: "cat",
			wantStdout: []byte(`this is some stdin text
and some more text

and some more again :P
`),
		},
		{
			name:       "stdin and stdout",
			stdin:      []byte("this is some stdin text"),
			command:    "sh",
			args:       []string{"-c", "echo 'hi there'; cat"},
			wantStdout: []byte("hi there\nthis is some stdin text"),
		},
		{
			name:    "stdin, stdout, and stderr",
			stdin:   []byte("this is some stdin text"),
			command: "sh",
			args: []string{
				"-c",
				`echo "hello world again"
echo "\n\noops broken\n\n" >&2
cat
`,
			},
			wantStdout: []byte("hello world again\nthis is some stdin text"),
			wantStderr: []byte("\n\noops broken\n\n\n"),
		},
		{
			name:    "error with no output",
			command: "sh",
			args:    []string{"-c", `exit 42`},
			wantErr: "exit status 42",
		},
		{
			name:    "error with stderr output",
			command: "sh",
			args: []string{
				"-c", `echo "\n\noops broken\n\n" >&2; exit 42`,
			},
			wantErr:    "exit status 42",
			wantStderr: []byte("\n\noops broken\n\n\n"),
		},
		{
			name:       "error with stdout output",
			command:    "sh",
			args:       []string{"-c", `echo 'hello world again'; exit 84`},
			wantStdout: []byte("hello world again\n"),
			wantErr:    "exit status 84",
		},
		{
			name:    "error with stdout and stderr output",
			command: "sh",
			args: []string{
				"-c",
				`echo "\n\noops broken\n\n" >&2
echo "hello world again"
exit 84
`,
			},
			wantStdout: []byte("hello world again\n"),
			wantErr:    "exit status 84",
			wantStderr: []byte("\n\noops broken\n\n\n"),
		},
		{
			name:    "error with dicarded stderr",
			command: "sh",
			args: []string{
				"-c",
				`echo "\n\noops broken\n\n" >&2
echo "hello world again"
exit 84
`,
			},
			discardStderr: true,
			wantStdout:    []byte("hello world again\n"),
			wantErr:       "exit status 84",
		},
		{
			name:    "error with discarded stdout",
			command: "sh",
			args: []string{
				"-c",
				`echo "\n\noops broken\n\n" >&2
echo "hello world again"
exit 84
`,
			},
			discardStdout: true,
			wantErr:       "exit status 84",
			wantStderr:    []byte("\n\noops broken\n\n\n"),
		},
		{
			name:    "error with discarded stdout and stderr",
			command: "sh",
			args: []string{
				"-c",
				`echo "\n\noops broken\n\n" >&2
echo "hello world again"
exit 84
`,
			},
			discardStdout: true,
			discardStderr: true,
			wantErr:       "exit status 84",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Local{}
			var stdin io.Reader
			if tt.stdin != nil {
				stdin = bytes.NewBuffer(tt.stdin)
			}
			var stdout io.ReadWriter
			if !tt.discardStdout {
				stdout = &bytes.Buffer{}
			}
			var stderr io.ReadWriter
			if !tt.discardStderr {
				stderr = &bytes.Buffer{}
			}

			err := r.Run(stdin, stdout, stderr, tt.command, tt.args...)

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}

			if !tt.discardStdout {
				if tt.wantStdout == nil {
					tt.wantStdout = []byte{}
				}
				b, err := io.ReadAll(stdout)
				require.NoError(t, err)
				assert.Equal(t, tt.wantStdout, b)
			}

			if stderr != nil {
				if tt.wantStderr == nil {
					tt.wantStderr = []byte{}
				}
				b, err := io.ReadAll(stderr)
				require.NoError(t, err)
				assert.Equal(t, tt.wantStderr, b)
			}
		})
	}
}
