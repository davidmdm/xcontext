package xcontext_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/davidmdm/xcontext"
	"github.com/stretchr/testify/require"
)

func TestWithCancelation(t *testing.T) {
	file, err := os.CreateTemp("", "acceptance-binary-*")
	require.NoError(t, err)

	coverDir := func() string {
		args := append([]string{}, os.Args...)
		slices.Reverse(args)
		for _, arg := range args {
			dir, ok := strings.CutPrefix(arg, "-test.gocoverdir=")
			if ok {
				return dir
			}
		}
		return ""
	}()

	t.Log("COVERDIR", coverDir)

	build := func() *exec.Cmd {
		if coverDir != "" {
			return CommandStandardIO("go", "build", "-coverpkg=./...", "-o", file.Name(), "./acceptance")
		}
		return CommandStandardIO("go", "build", "-o", file.Name(), "./acceptance")
	}()

	require.NoError(t, build.Run())
	require.NoError(t, file.Close())

	acceptanceCMD := func() (cmd *exec.Cmd, stdout *bytes.Buffer) {
		cmd = exec.Command(file.Name())
		if coverDir != "" {
			cmd.Env = append(cmd.Env, "GOCOVERDIR="+coverDir)
		}

		stdout = new(bytes.Buffer)
		cmd.Stdout = stdout

		return
	}

	t.Run("ctx canceled", func(t *testing.T) {
		acceptance, stdout := acceptanceCMD()
		acceptance.Env = append(acceptance.Env, "USE_CANCEL=true")

		require.NoError(t, acceptance.Run())

		require.Equal(t, "context canceled\n", stdout.String())
	})

	t.Run("dealine exceeded", func(t *testing.T) {
		acceptance, stdout := acceptanceCMD()
		acceptance.Env = append(acceptance.Env, "ROOT_TIMEOUT=10ms")

		require.NoError(t, acceptance.Run())

		require.Equal(t, "context deadline exceeded\n", stdout.String())
	})

	t.Run("sigint", func(t *testing.T) {
		acceptance, stdout := acceptanceCMD()

		require.NoError(t, acceptance.Start())

		// Give the process a chance to register the interceptors
		time.Sleep(50 * time.Millisecond)

		acceptance.Process.Signal(syscall.SIGINT)

		require.NoError(t, acceptance.Wait())
		require.Equal(t, "context canceled: received signal: interrupt\n", stdout.String())
	})
}

func TestSignalCause(t *testing.T) {
	cases := []struct {
		Name   string
		Err    error
		Signal os.Signal
	}{
		{
			Name:   "nil error",
			Err:    nil,
			Signal: nil,
		},
		{
			Name:   "error without signal",
			Err:    errors.New("some error"),
			Signal: nil,
		},
		{
			Name:   "exactly signal cause error",
			Err:    xcontext.SignalCancelError{Signal: syscall.SIGINT},
			Signal: syscall.SIGINT,
		},
		{
			Name:   "wrapped signal cause error",
			Err:    fmt.Errorf("wrapped: %w", xcontext.SignalCancelError{Signal: syscall.SIGINT}),
			Signal: syscall.SIGINT,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx, cancel := context.WithCancelCause(context.Background())
			cancel(tc.Err)

			require.Equal(t, tc.Signal, xcontext.SignalCause(ctx))
		})
	}
}

func CommandStandardIO(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(cmd.Env, os.Environ()...)

	return cmd
}
