// This file is adapted from various parts of the Docker CLI project: https://github.com/docker/cli/

package interactive

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	gosignal "os/signal"
	"runtime"
	"time"

	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/moby/sys/signal"
)

// resizeTtyTo resizes tty to specific height and width
func resizeTtyTo(
	ctx context.Context,
	logger *slog.Logger,
	docker client.ContainerAPIClient,
	id string,
	height uint,
	width uint,
	isExec bool,
) error {
	if height == 0 && width == 0 {
		return nil
	}

	options := container.ResizeOptions{
		Height: height,
		Width:  width,
	}

	var err error
	if isExec {
		err = docker.ContainerExecResize(ctx, id, options)
	} else {
		err = docker.ContainerResize(ctx, id, options)
	}

	if err != nil {
		logger.Debug(fmt.Sprintf("Error resize: %s\r", err))
	}
	return err
}

// resizeTty is to resize the tty with cli out's tty size
func resizeTty(
	ctx context.Context,
	logger *slog.Logger,
	docker client.ContainerAPIClient,
	streams command.Streams,
	id string,
	isExec bool,
) error {
	height, width := streams.Out().GetTtySize()
	return resizeTtyTo(ctx, logger, docker, id, height, width, isExec)
}

// initTtySize is to init the tty's size to the same as the window, if there is an error, it will retry 10 times.
func initTtySize(
	ctx context.Context,
	logger *slog.Logger,
	docker client.ContainerAPIClient,
	streams command.Streams,
	id string,
	isExec bool,
	resizeTtyFunc func(ctx context.Context, logger *slog.Logger, docker client.ContainerAPIClient, streams command.Streams, id string, isExec bool) error,
) {
	rttyFunc := resizeTtyFunc
	if rttyFunc == nil {
		rttyFunc = resizeTty
	}
	if err := rttyFunc(ctx, logger, docker, streams, id, isExec); err != nil {
		go func() {
			var err error
			for retry := 0; retry < 10; retry++ {
				time.Sleep(time.Duration(retry+1) * 10 * time.Millisecond)
				if err = rttyFunc(ctx, logger, docker, streams, id, isExec); err == nil {
					break
				}
			}
			if err != nil {
				fmt.Fprintln(streams.Err(), "failed to resize tty, using default size")
			}
		}()
	}
}

// MonitorTtySize updates the container tty size when the terminal tty changes size
func MonitorTtySize(
	ctx context.Context,
	logger *slog.Logger,
	docker client.ContainerAPIClient,
	streams command.Streams,
	id string,
	isExec bool,
) error {
	initTtySize(ctx, logger, docker, streams, id, isExec, resizeTty)
	if runtime.GOOS == "windows" {
		go func() {
			prevH, prevW := streams.Out().GetTtySize()
			for {
				time.Sleep(time.Millisecond * 250)
				h, w := streams.Out().GetTtySize()

				if prevW != w || prevH != h {
					resizeTty(ctx, logger, docker, streams, id, isExec)
				}
				prevH = h
				prevW = w
			}
		}()
	} else {
		sigchan := make(chan os.Signal, 1)
		gosignal.Notify(sigchan, signal.SIGWINCH)
		go func() {
			for range sigchan {
				resizeTty(ctx, logger, docker, streams, id, isExec)
			}
		}()
	}
	return nil
}
