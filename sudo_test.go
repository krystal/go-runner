package runner

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	mock_runner "github.com/krystal/go-runner/mock"
	"github.com/romdo/gomockctx"
	"github.com/stretchr/testify/assert"
)

func TestSudo_Run(t *testing.T) {
	type fields struct {
		User string
		Args []string
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
		fields      fields
		args        args
		err         error
		wantCommand string
		wantArgs    []string
		wantErr     string
	}{
		{
			name:   "sudo",
			fields: fields{},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "sudo",
			wantArgs:    []string{"-n", "--", "docker", "ps", "-a"},
		},
		{
			name:   "stdin",
			fields: fields{},
			args: args{
				stdin:   bytes.NewBufferString("foo\nbar"),
				stdout:  nil,
				stderr:  nil,
				command: "docker",
				args:    []string{"kill", "-s", "HUP"},
			},
			wantCommand: "sudo",
			wantArgs:    []string{"-n", "--", "docker", "kill", "-s", "HUP"},
		},
		{
			name:   "stdout",
			fields: fields{},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  nil,
				command: "docker",
				args:    []string{"stop", "foo"},
			},
			wantCommand: "sudo",
			wantArgs:    []string{"-n", "--", "docker", "stop", "foo"},
		},
		{
			name:   "stderr",
			fields: fields{},
			args: args{
				stdin:   nil,
				stdout:  nil,
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"stop", "foo"},
			},
			wantCommand: "sudo",
			wantArgs:    []string{"-n", "--", "docker", "stop", "foo"},
		},
		{
			name: "with User",
			fields: fields{
				User: "barfoo",
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "sudo",
			wantArgs: []string{
				"-n", "-u", "barfoo", "--", "docker", "ps", "-a",
			},
		},
		{
			name: "with Args",
			fields: fields{
				Args: []string{"-g", "other", "-d", "/opt/thing/data"},
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "sudo",
			wantArgs: []string{
				"-n", "-g", "other", "-d", "/opt/thing/data",
				"--", "docker", "ps", "-a",
			},
		},
		{
			name: "with User and Args",
			fields: fields{
				User: "barfoo",
				Args: []string{"-g", "other", "-d", "/opt/thing/data"},
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "sudo",
			wantArgs: []string{
				"-n", "-u", "barfoo", "-g", "other", "-d", "/opt/thing/data",
				"--", "docker", "ps", "-a",
			},
		},
		{
			name:   "error",
			fields: fields{},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "zfs",
				args:    []string{"list"},
			},
			err:         errors.New("zfs: command not found"),
			wantCommand: "sudo",
			wantArgs:    []string{"-n", "--", "zfs", "list"},
			wantErr:     "zfs: command not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			r := mock_runner.NewMockRunner(ctrl)
			r.EXPECT().Run(
				tt.args.stdin,
				tt.args.stdout,
				tt.args.stderr,
				tt.wantCommand,
				tt.wantArgs,
			).Return(tt.err)

			s := &Sudo{
				Runner: r,
				User:   tt.fields.User,
				Args:   tt.fields.Args,
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

func TestSudo_RunContext(t *testing.T) {
	ctx := gomockctx.New(context.Background())

	type fields struct {
		User string
		Args []string
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
		fields      fields
		args        args
		err         error
		wantCommand string
		wantArgs    []string
		wantErr     string
	}{
		{
			name:   "sudo",
			fields: fields{},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "sudo",
			wantArgs:    []string{"-n", "--", "docker", "ps", "-a"},
		},
		{
			name:   "stdin",
			fields: fields{},
			args: args{
				ctx:     ctx,
				stdin:   bytes.NewBufferString("foo\nbar"),
				stdout:  nil,
				stderr:  nil,
				command: "docker",
				args:    []string{"kill", "-s", "HUP"},
			},
			wantCommand: "sudo",
			wantArgs:    []string{"-n", "--", "docker", "kill", "-s", "HUP"},
		},
		{
			name:   "stdout",
			fields: fields{},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  nil,
				command: "docker",
				args:    []string{"stop", "foo"},
			},
			wantCommand: "sudo",
			wantArgs:    []string{"-n", "--", "docker", "stop", "foo"},
		},
		{
			name:   "stderr",
			fields: fields{},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  nil,
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"stop", "foo"},
			},
			wantCommand: "sudo",
			wantArgs:    []string{"-n", "--", "docker", "stop", "foo"},
		},
		{
			name: "with User",
			fields: fields{
				User: "barfoo",
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "sudo",
			wantArgs: []string{
				"-n", "-u", "barfoo", "--", "docker", "ps", "-a",
			},
		},
		{
			name: "with Args",
			fields: fields{
				Args: []string{"-g", "other", "-d", "/opt/thing/data"},
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "sudo",
			wantArgs: []string{
				"-n", "-g", "other", "-d", "/opt/thing/data",
				"--", "docker", "ps", "-a",
			},
		},
		{
			name: "with User and Args",
			fields: fields{
				User: "barfoo",
				Args: []string{"-g", "other", "-d", "/opt/thing/data"},
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "docker",
				args:    []string{"ps", "-a"},
			},
			wantCommand: "sudo",
			wantArgs: []string{
				"-n", "-u", "barfoo", "-g", "other", "-d", "/opt/thing/data",
				"--", "docker", "ps", "-a",
			},
		},
		{
			name:   "error",
			fields: fields{},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  &bytes.Buffer{},
				command: "zfs",
				args:    []string{"list"},
			},
			err:         errors.New("zfs: command not found"),
			wantCommand: "sudo",
			wantArgs:    []string{"-n", "--", "zfs", "list"},
			wantErr:     "zfs: command not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			r := mock_runner.NewMockRunner(ctrl)
			r.EXPECT().RunContext(
				gomockctx.Eq(tt.args.ctx),
				tt.args.stdin,
				tt.args.stdout,
				tt.args.stderr,
				tt.wantCommand,
				tt.wantArgs,
			).Return(tt.err)

			s := &Sudo{
				Runner: r,
				User:   tt.fields.User,
				Args:   tt.fields.Args,
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

func TestSudo_Env(t *testing.T) {
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
			r.EXPECT().Env(tt.args.env)

			s := &Sudo{Runner: r}

			s.Env(tt.args.env...)
		})
	}
}
