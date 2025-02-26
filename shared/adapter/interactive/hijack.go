// This file is adapted from various parts of the Docker CLI project: https://github.com/docker/cli/

package interactive

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"sync"

	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/term"
)

// The default escape key sequence: ctrl-p, ctrl-q
var defaultEscapeKeys = []byte{16, 17}

// A HijackedIOStreamer handles copying input to and output from streams to the
// connection.
type HijackedIOStreamer struct {
	logger       *slog.Logger
	streams      command.Streams
	inputStream  io.ReadCloser
	outputStream io.Writer
	errorStream  io.Writer

	tty  bool
	resp types.HijackedResponse
}

// Creates a new HijackedIOStreamer.
func NewHijackedIOStreamer(tty bool, logger *slog.Logger, streams command.Streams, resp types.HijackedResponse) *HijackedIOStreamer {
	return &HijackedIOStreamer{
		logger:       logger,
		streams:      streams,
		inputStream:  streams.In(),
		outputStream: streams.Out(),
		errorStream:  streams.Err(),
		tty:          tty,
		resp:         resp,
	}
}

// Stream handles setting up the IO and then begins streaming stdin/stdout
// to/from the hijacked connection, blocking until it is either done reading
// output, the user inputs the detach key sequence when in TTY mode, or when
// the given context is cancelled.
func (h *HijackedIOStreamer) Stream(ctx context.Context) error {
	restoreInput, err := h.setupInput()
	if err != nil {
		return fmt.Errorf("unable to setup input stream: %s", err)
	}

	defer restoreInput()

	outputDone := h.beginOutputStream(restoreInput)
	inputDone, detached := h.beginInputStream(restoreInput)

	select {
	case err := <-outputDone:
		return err
	case <-inputDone:
		// Input stream has closed.
		if h.outputStream != nil || h.errorStream != nil {
			// Wait for output to complete streaming.
			select {
			case err := <-outputDone:
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	case err := <-detached:
		// Got a detach key sequence.
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *HijackedIOStreamer) setupInput() (restore func(), err error) {
	if h.inputStream == nil || !h.tty {
		// No need to setup input TTY.
		// The restore func is a nop.
		return func() {}, nil
	}

	if err := setRawTerminal(h.streams); err != nil {
		return nil, fmt.Errorf("unable to set IO streams as raw terminal: %s", err)
	}

	// Use sync.Once so we may call restore multiple times but ensure we
	// only restore the terminal once.
	var restoreOnce sync.Once
	restore = func() {
		restoreOnce.Do(func() {
			_ = restoreTerminal(h.streams, h.inputStream)
		})
	}

	// Wrap the input to detect detach escape sequence.
	escapeKeys := defaultEscapeKeys
	h.inputStream = ioutils.NewReadCloserWrapper(term.NewEscapeProxy(h.inputStream, escapeKeys), h.inputStream.Close)

	return restore, nil
}

func (h *HijackedIOStreamer) beginOutputStream(restoreInput func()) <-chan error {
	if h.outputStream == nil && h.errorStream == nil {
		// There is no need to copy output.
		return nil
	}

	outputDone := make(chan error)
	go func() {
		var err error

		// When TTY is ON, use regular copy
		if h.outputStream != nil && h.tty {
			_, err = io.Copy(h.outputStream, h.resp.Reader)
			// We should restore the terminal as soon as possible
			// once the connection ends so any following print
			// messages will be in normal type.
			restoreInput()
		} else {
			_, err = stdcopy.StdCopy(h.outputStream, h.errorStream, h.resp.Reader)
		}

		h.logger.Debug("[hijack] End of stdout")

		if err != nil {
			h.logger.Debug(fmt.Sprintf("Error receiveStdout: %s", err))
		}

		outputDone <- err
	}()

	return outputDone
}

func (h *HijackedIOStreamer) beginInputStream(restoreInput func()) (doneC <-chan struct{}, detachedC <-chan error) {
	inputDone := make(chan struct{})
	detached := make(chan error)

	go func() {
		if h.inputStream != nil {
			_, err := io.Copy(h.resp.Conn, h.inputStream)
			// We should restore the terminal as soon as possible
			// once the connection ends so any following print
			// messages will be in normal type.
			restoreInput()

			h.logger.Debug("[hijack] End of stdin")

			if _, ok := err.(term.EscapeError); ok {
				detached <- err
				return
			}

			if err != nil {
				// This error will also occur on the receive
				// side (from stdout) where it will be
				// propagated back to the caller.
				h.logger.Debug(fmt.Sprintf("Error sendStdin: %s", err))
			}
		}

		if err := h.resp.CloseWrite(); err != nil {
			h.logger.Debug(fmt.Sprintf("Couldn't send EOF: %s", err))
		}

		close(inputDone)
	}()

	return inputDone, detached
}

func setRawTerminal(streams command.Streams) error {
	if err := streams.In().SetRawTerminal(); err != nil {
		return err
	}
	return streams.Out().SetRawTerminal()
}

func restoreTerminal(streams command.Streams, in io.Closer) error {
	streams.In().RestoreTerminal()
	streams.Out().RestoreTerminal()
	// WARNING: DO NOT REMOVE THE OS CHECKS !!!
	// For some reason this Close call blocks on darwin..
	// As the client exits right after, simply discard the close
	// until we find a better solution.
	//
	// This can also cause the client on Windows to get stuck in Win32 CloseHandle()
	// in some cases. See https://github.com/docker/docker/issues/28267#issuecomment-288237442
	// Tracked internally at Microsoft by VSO #11352156. In the
	// Windows case, you hit this if you are using the native/v2 console,
	// not the "legacy" console, and you start the client in a new window. eg
	// `start docker run --rm -it microsoft/nanoserver cmd /s /c echo foobar`
	// will hang. Remove start, and it won't repro.
	if in != nil && runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
		return in.Close()
	}
	return nil
}
