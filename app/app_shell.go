package app

import (
	"os/exec"

	"github.com/go-cmd/cmd"
)

var shell string

func init() {
	for _, n := range []string{"bash", "sh"} {
		path, err := exec.LookPath(n)
		if err == nil {
			shell = path
			break
		}
	}
}

type ShellOptions struct {
	Out func(string, string)
	Err func(string, string)
}

func Shell(command string, options ShellOptions) (cmd.Status, error) {
	if options.Out == nil {
		options.Out = func(string, string) {}
	}
	if options.Err == nil {
		options.Err = func(string, string) {}
	}

	// see: https://github.com/go-cmd/cmd/blob/master/examples/blocking-streaming/main.go
	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	// Create Cmd with options
	envCmd := cmd.NewCmdOptions(cmdOptions, shell, "-c", command)

	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for envCmd.Stdout != nil || envCmd.Stderr != nil {
			select {
			case line, open := <-envCmd.Stdout:
				if !open {
					envCmd.Stdout = nil
					continue
				}
				options.Out("CommandJob", line)
			case line, open := <-envCmd.Stderr:
				if !open {
					envCmd.Stderr = nil
					continue
				}
				options.Err("CommandJob", line)
			}
		}
	}()

	// Run and wait for Cmd to return
	status := <-envCmd.Start()

	// Wait for goroutine to print everything
	<-doneChan

	return status, nil
}
