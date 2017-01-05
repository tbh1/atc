package container

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"code.cloudfoundry.org/garden"
	"github.com/tedsuo/ifrit"
)

type Wrapper struct {
	Container garden.Container
}

var ErrAborted = errors.New("script aborted")

type ErrScriptFailed struct {
	Path       string
	Args       []string
	ExitStatus int

	Stderr string
}

func (err ErrScriptFailed) Error() string {
	msg := fmt.Sprintf(
		"script '%s %v' failed: exit status %d",
		err.Path,
		err.Args,
		err.ExitStatus,
	)

	if len(err.Stderr) > 0 {
		msg += "\n\nstderr:\n" + err.Stderr
	}

	return msg
}

func (cw *Wrapper) RunScript(
	path string,
	args []string,
	input interface{},
	output interface{},
) ifrit.Runner {
	return ifrit.RunFunc(func(signals <-chan os.Signal, ready chan<- struct{}) error {
		request, err := json.Marshal(input)
		if err != nil {
			return err
		}

		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)

		processIO := garden.ProcessIO{
			Stdin:  bytes.NewBuffer(request),
			Stdout: stdout,
			Stderr: stderr,
		}

		var process garden.Process

		process, err = cw.Container.Run(garden.ProcessSpec{
			Path: path,
			Args: args,
		}, processIO)
		if err != nil {
			return err
		}
		close(ready)

		processExited := make(chan struct{})

		var processStatus int
		var processErr error

		go func() {
			processStatus, processErr = process.Wait()
			close(processExited)
		}()

		select {
		case <-processExited:
			if processErr != nil {
				return processErr
			}

			if processStatus != 0 {
				return ErrScriptFailed{
					Path:       path,
					Args:       args,
					ExitStatus: processStatus,
					Stderr:     stderr.String(),
				}
			}

			return json.Unmarshal(stdout.Bytes(), output)

		case <-signals:
			cw.Container.Stop(false)
			<-processExited
			return ErrAborted
		}
	})

}
