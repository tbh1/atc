package resource

import (
	"bytes"
	"encoding/json"

	"code.cloudfoundry.org/garden"
	"github.com/concourse/atc"
	"github.com/concourse/atc/cessna"
	"os"
)

func NewPutCommandProcess(container garden.Container, resource Resource, params atc.Params, artifactsDirectory string) (*putCommandProcess, error) {
	p := &cessna.ContainerProcess{
		Container: container,
		ProcessSpec: garden.ProcessSpec{
			Path: "/opt/resource/out",
			Args: []string{artifactsDirectory},
		},
	}

	i := OutRequest{
		Source: resource.Source,
		Params: params,
	}

	input, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	p.ProcessIO.Stdin = bytes.NewBuffer(input)
	p.ProcessIO.Stdout = &stdout
	p.ProcessIO.Stderr = &stderr

	return &putCommandProcess{
		ContainerProcess: p,
		out:              &stdout,
		err:              &stderr,
	}, nil
}

type putCommandProcess struct {
	*cessna.ContainerProcess

	out *bytes.Buffer
	err *bytes.Buffer
}

func (g *putCommandProcess) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	err := g.ContainerProcess.Run(signals, ready)
	if err != nil {
		switch e := err.(type) {
		case cessna.ErrScriptFailed:
			e.Stderr = string(g.err.Bytes())

			err = e
		}
	}

	return err
}

func (g *putCommandProcess) Response() (OutResponse, error) {
	var r OutResponse

	err := json.NewDecoder(g.out).Decode(&r)
	if err != nil {
		return OutResponse{}, err
	}

	return r, nil
}
