package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	mock_runner "github.com/krystal/go-runner/mock"
	"github.com/romdo/gomockctx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type fakeTestingT struct {
	Messages []string
}

func (f *fakeTestingT) Logf(format string, args ...interface{}) {
	if f == nil {
		return
	}

	f.Messages = append(f.Messages, fmt.Sprintf(format, args...))
}

func TestTesting_Run(t *testing.T) {
	type fields struct {
		T *fakeTestingT
	}
	type args struct {
		stdin   io.Reader
		stdout  io.Writer
		stderr  io.Writer
		command string
		args    []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		err     error
		wantErr string
		wantLog []string
	}{
		{
			name: "no T",
			args: args{
				stdin:   nil,
				stdout:  nil,
				stderr:  nil,
				command: "echo",
				args:    []string{"-n", "hello world"},
			},
			wantLog: []string{},
		},
		{
			name: "echo",
			fields: fields{
				T: &fakeTestingT{},
			},
			args: args{
				stdin:   nil,
				stdout:  nil,
				stderr:  nil,
				command: "echo",
				args:    []string{"-n", "hello world"},
			},
			wantLog: []string{
				`runner.Run: command=echo args=["-n","hello world"]`,
			},
		},
		{
			name: "stdin",
			fields: fields{
				T: &fakeTestingT{},
			},
			args: args{
				stdin:   bytes.NewBufferString("foo\nbar"),
				stdout:  nil,
				stderr:  nil,
				command: "echo",
				args:    []string{"hi", "john"},
			},
			wantLog: []string{
				`runner.Run: command=echo args=["hi","john"]`,
			},
		},
		{
			name: "stdout",
			fields: fields{
				T: &fakeTestingT{},
			},
			args: args{
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  nil,
				command: "echo",
				args:    []string{"hi", "jane"},
			},
			wantLog: []string{
				`runner.Run: command=echo args=["hi","jane"]`,
			},
		},
		{
			name: "stderr",
			fields: fields{
				T: &fakeTestingT{},
			},
			args: args{
				stdin:   nil,
				stdout:  nil,
				stderr:  &bytes.Buffer{},
				command: "ps",
				args:    []string{"-a", "-ux"},
			},
			wantLog: []string{
				`runner.Run: command=ps args=["-a","-ux"]`,
			},
		},
		{
			name: "error",
			fields: fields{
				T: &fakeTestingT{},
			},
			args: args{
				stdin:   nil,
				stdout:  nil,
				stderr:  &bytes.Buffer{},
				command: "false",
				args:    []string{},
			},
			err:     errors.New("exit status 1"),
			wantErr: "exit status 1",
			wantLog: []string{
				`runner.Run: command=false args=[]`,
			},
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
				tt.args.command,
				tt.args.args,
			).Return(tt.err)

			tr := &Testing{
				Runner:   r,
				TestingT: tt.fields.T,
			}

			err := tr.Run(
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

			if tt.fields.T != nil {
				assert.Equal(t, tt.wantLog, tt.fields.T.Messages)
			} else {
				assert.Empty(t, tt.wantLog)
			}
		})
	}
}

func TestTesting_RunContext(t *testing.T) {
	ctx := gomockctx.New(context.Background())

	type fields struct {
		T *fakeTestingT
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
		name    string
		fields  fields
		args    args
		err     error
		wantErr string
		wantLog []string
	}{
		{
			name: "no T",
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  nil,
				stderr:  nil,
				command: "echo",
				args:    []string{"-n", "hello world"},
			},
			wantLog: []string{},
		},
		{
			name: "echo",
			fields: fields{
				T: &fakeTestingT{},
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  nil,
				stderr:  nil,
				command: "echo",
				args:    []string{"-n", "hello world"},
			},
			wantLog: []string{
				`runner.RunContext: command=echo args=["-n","hello world"]`,
			},
		},
		{
			name: "stdin",
			fields: fields{
				T: &fakeTestingT{},
			},
			args: args{
				ctx:     ctx,
				stdin:   bytes.NewBufferString("foo\nbar"),
				stdout:  nil,
				stderr:  nil,
				command: "echo",
				args:    []string{"hi", "john"},
			},
			wantLog: []string{
				`runner.RunContext: command=echo args=["hi","john"]`,
			},
		},
		{
			name: "stdout",
			fields: fields{
				T: &fakeTestingT{},
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  &bytes.Buffer{},
				stderr:  nil,
				command: "echo",
				args:    []string{"hi", "jane"},
			},
			wantLog: []string{
				`runner.RunContext: command=echo args=["hi","jane"]`,
			},
		},
		{
			name: "stderr",
			fields: fields{
				T: &fakeTestingT{},
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  nil,
				stderr:  &bytes.Buffer{},
				command: "ps",
				args:    []string{"-a", "-ux"},
			},
			wantLog: []string{
				`runner.RunContext: command=ps args=["-a","-ux"]`,
			},
		},
		{
			name: "error",
			fields: fields{
				T: &fakeTestingT{},
			},
			args: args{
				ctx:     ctx,
				stdin:   nil,
				stdout:  nil,
				stderr:  &bytes.Buffer{},
				command: "false",
				args:    []string{},
			},
			err:     errors.New("exit status 1"),
			wantErr: "exit status 1",
			wantLog: []string{
				`runner.RunContext: command=false args=[]`,
			},
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
				tt.args.command,
				tt.args.args,
			).Return(tt.err)

			tr := &Testing{
				Runner:   r,
				TestingT: tt.fields.T,
			}

			err := tr.RunContext(
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

			if tt.fields.T != nil {
				assert.Equal(t, tt.wantLog, tt.fields.T.Messages)
			} else {
				assert.Empty(t, tt.wantLog)
			}
		})
	}
}

func TestTesting_Env(t *testing.T) {
	type fields struct {
		T      *fakeTestingT
		LogEnv bool
	}
	type args struct {
		env []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLog []string
	}{
		{
			name: "empty and no LogEnv",
			fields: fields{
				T:      &fakeTestingT{},
				LogEnv: false,
			},
			args: args{
				env: []string{},
			},
		},
		{
			name: "empty and LogEnv",
			fields: fields{
				T:      &fakeTestingT{},
				LogEnv: true,
			},
			args: args{
				env: []string{},
			},
			wantLog: []string{
				"runner.Env: vars=[]",
			},
		},
		{
			name: "one var",
			fields: fields{
				T:      &fakeTestingT{},
				LogEnv: false,
			},
			args: args{
				env: []string{"foo=bar"},
			},
		},
		{
			name: "one var and LogEnv",
			fields: fields{
				T:      &fakeTestingT{},
				LogEnv: true,
			},
			args: args{
				env: []string{"foo=bar"},
			},
			wantLog: []string{
				`runner.Env: vars=["foo=bar"]`,
			},
		},
		{
			name: "many vars",
			fields: fields{
				T:      &fakeTestingT{},
				LogEnv: false,
			},
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
		{
			name: "many vars and LogEnv",
			fields: fields{
				T:      &fakeTestingT{},
				LogEnv: true,
			},
			args: args{
				env: []string{
					"foo=bar",
					"foo=bar",
					"foz=baz",
					"nope=why",
					"hello=world",
				},
			},
			wantLog: []string{
				`runner.Env: vars=[` +
					`"foo=bar",` +
					`"foo=bar",` +
					`"foz=baz",` +
					`"nope=why",` +
					`"hello=world"` +
					`]`,
			},
		},
		{
			name: "no T",
			fields: fields{
				LogEnv: true,
			},
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

			tr := &Testing{
				Runner:   r,
				TestingT: tt.fields.T,
				LogEnv:   tt.fields.LogEnv,
			}

			tr.Env(tt.args.env...)

			if tt.fields.T != nil {
				assert.Equal(t, tt.wantLog, tt.fields.T.Messages)
			} else {
				assert.Empty(t, tt.wantLog)
			}
		})
	}
}
