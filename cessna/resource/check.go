package resource

import (
	"bytes"
	"encoding/json"

	"code.cloudfoundry.org/garden"
	"github.com/concourse/atc"
	"github.com/concourse/atc/cessna"
)

func NewCheckCommandProcess(container garden.Container, resource Resource, version *atc.Version) (*checkCommandProcess, error) {
	p := &cessna.ContainerProcess{
		Container: container,
		ProcessSpec: garden.ProcessSpec{
			Path: "/opt/resource/check",
		},
	}

	i := CheckRequest{
		Source:  resource.Source,
	}

	if version != nil {
		i.Version = *version
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

	return &checkCommandProcess{
		ContainerProcess: p,
		out:              &stdout,
		err:              &stderr,
	}, nil
}

type checkCommandProcess struct {
	*cessna.ContainerProcess

	out *bytes.Buffer
	err *bytes.Buffer
}

func (c *checkCommandProcess) Response() (CheckResponse, error) {
	var o CheckResponse

	err := json.NewDecoder(c.out).Decode(&o)
	if err != nil {
		return nil, err
	}

	return o, nil
}
