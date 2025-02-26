package interactive

import (
	"github.com/docker/cli/cli/streams"
	"github.com/moby/term"
)

// A wrapper for stdin, stdout, and stderr streams, compliant with Docker's command.Streams interface
type StandardStreamWrapper struct {
	// The standard input stream
	stdin *streams.In

	// The standard output stream
	stdout *streams.Out

	// The standard error stream
	stderr *streams.Out
}

// Create a new standard stream wrapper
func NewStandardStreamWrapper(ttyMode bool) *StandardStreamWrapper {
	stdin, stdout, stderr := term.StdStreams()
	wrapper := &StandardStreamWrapper{
		stdin:  streams.NewIn(stdin),
		stdout: streams.NewOut(stdout),
	}
	// For TTY mode, use stdout for stderr
	if ttyMode {
		wrapper.stderr = wrapper.stdout
	} else {
		wrapper.stderr = streams.NewOut(stderr)
	}
	return wrapper
}

// Get the standard input stream
func (s *StandardStreamWrapper) In() *streams.In {
	return s.stdin
}

// Get the standard output stream
func (s *StandardStreamWrapper) Out() *streams.Out {
	return s.stdout
}

// Get the standard error stream
func (s *StandardStreamWrapper) Err() *streams.Out {
	return s.stderr
}
