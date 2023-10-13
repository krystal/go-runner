package runner

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	mock_runner "github.com/krystal/go-runner/mock"
	"github.com/romdo/gomockctx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSSHCLI_Run(t *testing.T) {
	type fields struct {
		Destination  string
		Port         int
		IdentityFile string
		Login        string
		Args         []string
	}
	type args struct {
		stdin   io.Reader
		stdout  io.Writer
		stderr  io.Writer
		command string
		args    []string
	}
	tests := []struct {
		name        string
		env         []string
		fields      fields
		args        args
		err         error
		wantCommand string
		wantArgs    []string
		wantErr     string
	}{
		{
			name: "ssh hostname",
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"narnia.local", "--", "docker", "ps", "-a",
			},
		},
		{
			name: "ssh user@hostname",
			fields: fields{
				Destination: "darrin@narnia.local",
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"darrin@narnia.local", "--", "docker", "ps", "-a",
			},
		},
		{
			name: "ssh ssh://user@hostname:port",
			fields: fields{
				Destination: "ssh://darrin@narnia.local:322",
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"ssh://darrin@narnia.local:322",
				"--", "docker", "ps", "-a",
			},
		},
		{
			name: "stdin",
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				stdin:   bytes.NewBufferString("foo\nbar"),
				stdout:  nil,
				stderr:  nil,
				command: "docker",
				args:    []string{"kill", "-s", "HUP"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"narnia.local", "--", "docker", "kill", "-s", "HUP",
			},
		},
		{
			name: "stdout",
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  nil,
				command: "docker",
				args:    []string{"stop", "foo"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"narnia.local", "--", "docker", "stop", "foo",
			},
		},
		{
			name: "stderr",
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				stdin:   nil,
				stdout:  nil,
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"stop", "foo"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"narnia.local", "--", "docker", "stop", "foo",
			},
		},
		{
			name: "with Port",
			fields: fields{
				Destination: "narnia.local",
				Port:        322,
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"-p", "322", "narnia.local", "--", "docker", "ps", "-a",
			},
		},
		{
			name: "with IdentityFile",
			fields: fields{
				Destination:  "narnia.local",
				IdentityFile: "/home/darrin/.ssh/id_other",
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"-i", "/home/darrin/.ssh/id_other", "narnia.local",
				"--", "docker", "ps", "-a",
			},
		},
		{
			name: "with Login",
			fields: fields{
				Destination: "narnia.local",
				Login:       "barfoo",
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"-l", "barfoo", "narnia.local", "--", "docker", "ps", "-a",
			},
		},
		{
			name: "with Env",
			env:  []string{"FOO=BAR", "PORT=8080"},
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "myapp",
				args:    []string{"run", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"narnia.local",
				"--", "env", "FOO=BAR", "PORT=8080", "myapp", "run", "-a",
			},
		},
		{
			name: "with Args",
			fields: fields{
				Destination: "narnia.local",
				Args:        []string{"-C", "-o", "AddKeysToAgent=yes"},
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"-C", "-o", "AddKeysToAgent=yes", "narnia.local",
				"--", "docker", "ps", "-a",
			},
		},
		{
			name: "with Port, IdentityFile, Login, Args and Env",
			env:  []string{"FOO=BAR", "PORT=8080"},
			fields: fields{
				Destination:  "narnia.local",
				Port:         322,
				IdentityFile: "/home/darrin/.ssh/id_other",
				Login:        "barfoo",
				Args:         []string{"-C", "-o", "AddKeysToAgent=yes"},
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"-p", "322",
				"-i", "/home/darrin/.ssh/id_other",
				"-l", "barfoo",
				"-C", "-o", "AddKeysToAgent=yes",
				"narnia.local",
				"--", "env", "FOO=BAR", "PORT=8080", "docker", "ps", "-a",
			},
		},
		{
			name:   "no destination",
			fields: fields{},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "zfs",
				args:    []string{"list"},
			},
			wantErr: ErrSSHCLINoDestination.Error(),
		},
		{
			name: "error",
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "zfs",
				args:    []string{"list"},
			},
			err:         errors.New("zfs: command not found"),
			wantCommand: "ssh",
			wantArgs:    []string{"narnia.local", "--", "zfs", "list"},
			wantErr:     "zfs: command not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			r := mock_runner.NewMockRunner(ctrl)
			if tt.wantCommand != "" {
				r.EXPECT().Run(
					tt.args.stdin,
					tt.args.stdout,
					tt.args.stderr,
					tt.wantCommand,
					tt.wantArgs,
				).Return(tt.err)
			}

			s := &SSHCLI{
				Runner:       r,
				Destination:  tt.fields.Destination,
				Port:         tt.fields.Port,
				IdentityFile: tt.fields.IdentityFile,
				Login:        tt.fields.Login,
				Args:         tt.fields.Args,
			}

			if len(tt.env) > 0 {
				s.Env(tt.env...)
			}

			err := s.Run(
				tt.args.stdin,
				tt.args.stdout,
				tt.args.stderr,
				tt.args.command,
				tt.args.args...,
			)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSSHCLI_RunContext(t *testing.T) {
	ctx := gomockctx.New(context.Background())

	type fields struct {
		Destination  string
		Port         int
		IdentityFile string
		Login        string
		Args         []string
	}

	type args struct {
		ctx     context.Context
		stdin   io.Reader
		stdout  io.Writer
		stderr  io.Writer
		command string
		args    []string
	}
	tests := []struct {
		name        string
		env         []string
		fields      fields
		args        args
		err         error
		wantCommand string
		wantArgs    []string
		wantErr     string
	}{
		{
			name: "ssh hostname",
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"narnia.local", "--", "docker", "ps", "-a",
			},
		},
		{
			name: "ssh user@hostname",
			fields: fields{
				Destination: "darrin@narnia.local",
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"darrin@narnia.local", "--", "docker", "ps", "-a",
			},
		},
		{
			name: "ssh ssh://user@hostname:port",
			fields: fields{
				Destination: "ssh://darrin@narnia.local:322",
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"ssh://darrin@narnia.local:322",
				"--", "docker", "ps", "-a",
			},
		},
		{
			name: "stdin",
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				ctx:     ctx,
				stdin:   bytes.NewBufferString("foo\nbar"),
				stdout:  nil,
				stderr:  nil,
				command: "docker",
				args:    []string{"kill", "-s", "HUP"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"narnia.local", "--", "docker", "kill", "-s", "HUP",
			},
		},
		{
			name: "stdout",
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  nil,
				command: "docker",
				args:    []string{"stop", "foo"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"narnia.local", "--", "docker", "stop", "foo",
			},
		},
		{
			name: "stderr",
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  nil,
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"stop", "foo"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"narnia.local", "--", "docker", "stop", "foo",
			},
		},
		{
			name: "with Port",
			fields: fields{
				Destination: "narnia.local",
				Port:        322,
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"-p", "322", "narnia.local", "--", "docker", "ps", "-a",
			},
		},
		{
			name: "with IdentityFile",
			fields: fields{
				Destination:  "narnia.local",
				IdentityFile: "/home/darrin/.ssh/id_other",
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"-i", "/home/darrin/.ssh/id_other", "narnia.local",
				"--", "docker", "ps", "-a",
			},
		},
		{
			name: "with Login",
			fields: fields{
				Destination: "narnia.local",
				Login:       "barfoo",
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"-l", "barfoo", "narnia.local", "--", "docker", "ps", "-a",
			},
		},
		{
			name: "with Env",
			env:  []string{"FOO=BAR", "PORT=8080"},
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "myapp",
				args:    []string{"run", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"narnia.local",
				"--", "env", "FOO=BAR", "PORT=8080", "myapp", "run", "-a",
			},
		},
		{
			name: "with Args",
			fields: fields{
				Destination: "narnia.local",
				Args:        []string{"-C", "-o", "AddKeysToAgent=yes"},
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"-C", "-o", "AddKeysToAgent=yes", "narnia.local",
				"--", "docker", "ps", "-a",
			},
		},
		{
			name: "with Port, IdentityFile, Login, Args and Env",
			env:  []string{"FOO=BAR", "PORT=8080"},
			fields: fields{
				Destination:  "narnia.local",
				Port:         322,
				IdentityFile: "/home/darrin/.ssh/id_other",
				Login:        "barfoo",
				Args:         []string{"-C", "-o", "AddKeysToAgent=yes"},
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "ssh",
			wantArgs: []string{
				"-p", "322",
				"-i", "/home/darrin/.ssh/id_other",
				"-l", "barfoo",
				"-C", "-o", "AddKeysToAgent=yes",
				"narnia.local",
				"--", "env", "FOO=BAR", "PORT=8080", "docker", "ps", "-a",
			},
		},
		{
			name:   "no destination",
			fields: fields{},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "zfs",
				args:    []string{"list"},
			},
			wantErr: ErrSSHCLINoDestination.Error(),
		},
		{
			name: "error",
			fields: fields{
				Destination: "narnia.local",
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "zfs",
				args:    []string{"list"},
			},
			err:         errors.New("zfs: command not found"),
			wantCommand: "ssh",
			wantArgs:    []string{"narnia.local", "--", "zfs", "list"},
			wantErr:     "zfs: command not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			r := mock_runner.NewMockRunner(ctrl)
			if tt.wantCommand != "" {
				r.EXPECT().RunContext(
					gomockctx.Eq(tt.args.ctx),
					tt.args.stdin,
					tt.args.stdout,
					tt.args.stderr,
					tt.wantCommand,
					tt.wantArgs,
				).Return(tt.err)
			}

			s := &SSHCLI{
				Runner:       r,
				Destination:  tt.fields.Destination,
				Port:         tt.fields.Port,
				IdentityFile: tt.fields.IdentityFile,
				Login:        tt.fields.Login,
				Args:         tt.fields.Args,
			}

			if len(tt.env) > 0 {
				s.Env(tt.env...)
			}

			err := s.RunContext(
				tt.args.ctx,
				tt.args.stdin,
				tt.args.stdout,
				tt.args.stderr,
				tt.args.command,
				tt.args.args...,
			)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSSHCLI_Env(t *testing.T) {
	type args struct {
		env []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "empty",
			args: args{
				env: []string{},
			},
		},
		{
			name: "one var",
			args: args{
				env: []string{
					"foo=bar",
				},
			},
		},
		{
			name: "many vars",
			args: args{
				env: []string{
					"foo=bar",
					"foo=bar",
					"foz=baz",
					"nope=why",
					"hello=world",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			r := mock_runner.NewMockRunner(ctrl)

			s := &SSHCLI{Runner: r}
			s.Env(tt.args.env...)

			assert.Equal(t, tt.args.env, s.env)
		})
	}
}
