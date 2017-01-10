package resource

import (
	"bytes"
	"encoding/json"

	"code.cloudfoundry.org/garden"
	"github.com/concourse/atc"
	"github.com/concourse/atc/cessna"
)

func NewGetCommandProcess(container garden.Container, mount garden.BindMount, resource Resource, version *atc.Version, params atc.Params) (*getCommandProcess, error) {
	p := &cessna.ContainerProcess{
		Container: container,
		ProcessSpec: garden.ProcessSpec{
			Path: "/opt/resource/in",
			Args: []string{mount.DstPath},
		},
	}

	i := InRequest{
		Source:  resource.Source,
		Params:  params,
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

	return &getCommandProcess{
		ContainerProcess: p,
		out:              &stdout,
		err:              &stderr,
	}, nil
}

type getCommandProcess struct {
	*cessna.ContainerProcess

	out *bytes.Buffer
	err *bytes.Buffer
}

func (g *getCommandProcess) Response() (InResponse, error) {
	var r InResponse

	err := json.NewDecoder(g.out).Decode(&r)
	if err != nil {
		return InResponse{}, err
	}

	return r, nil
}
