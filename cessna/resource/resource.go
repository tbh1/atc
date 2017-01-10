package resource

import (
	"bytes"
	"encoding/json"
	"os"

	"io"

	"fmt"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/cessna"
	"github.com/concourse/atc/cessna/container"
	"github.com/concourse/baggageclaim"
	"github.com/tedsuo/ifrit"
)

type ResourceType struct {
	RootFSPath string
	Name       string
}

type Resource struct {
	ResourceType ResourceType
	Source       atc.Source
}

type resourceContainer struct {
	container.Wrapper
	resource Resource
}

func NewResource(resourceType ResourceType, source atc.Source) Resource {
	return Resource{
		ResourceType: resourceType,
		Source:       source,
	}
}

type ResourceManager struct {
	Logger lager.Logger
	Worker cessna.Worker
}

func NewResourceManagerFor(worker cessna.Worker) *ResourceManager {
	return &ResourceManager{
		Logger: lager.NewLogger("resourcemanager"),
		Worker: worker,
	}
}

func (r *ResourceManager) Check(resource Resource, version *atc.Version) ([]atc.Version, error) {
	// Import RootFS into Volume
	spec := baggageclaim.VolumeSpec{
		Strategy: baggageclaim.ImportStrategy{
			Path: resource.ResourceType.RootFSPath,
		},
		Privileged: true,
	}
	handle := "foobar"

	fmt.Println("HERE")
	parentVolume, err := r.Worker.BaggageClaimClient().CreateVolume(r.Logger.Session("create-parent-volume"), handle, spec)
	if err != nil {
		return nil, err
	}

	// COW of RootFS Volume
	spec = baggageclaim.VolumeSpec{
		Strategy: baggageclaim.COWStrategy{
			Parent: parentVolume,
		},
		Privileged: false,
	}
	rootFSVolume, err := r.Worker.BaggageClaimClient().CreateVolume(r.Logger.Session("create-cow-rootfs-volume"), "foobar2", spec)
	if err != nil {
		return nil, err
	}

	// Turn RootFS COW into Container
	gardenSpec := garden.ContainerSpec{
		Privileged: false,
		RootFSPath: rootFSVolume.Path(),
	}

	container, err := r.Worker.GardenClient().Create(gardenSpec)
	if err != nil {
		r.Logger.Error("failed-to-create-container-in-garden", err)
		return nil, err
	}

	// Run Check Command in Garden Container
	runner, err := NewCheckCommandRunner(resource, container)
	if err != nil {
		return nil, err
	}

	checking := ifrit.Invoke(runner)

	err = <-checking.Wait()
	if err != nil {
		return nil, err
	}

	return runner.Versions()
}

func (r *ResourceManager) Get(resource Resource, version *atc.Version, params atc.Params) (cessna.Volume, error) {
	//volume, err := r.Volumizer.COWFromRootFS(resource.ResourceType.RootFSPath)
	//
	//volumeForGet, err := r.Volumizer.CreateEmptyVolume("123", false)
	//if err != nil {
	//	return nil, err
	//}
	//
	//bindMounts := &[]garden.BindMount{{
	//		SrcPath: volumeForGet.Path(),
	//		DstPath: "/tmp/resource/get",
	//		Mode:    garden.BindMountModeRW,
	//	}}
	//
	//resourceContainer, err := r.Containerizer.CreateContainer(volume.Path(), bindMounts)
	//if err != nil {
	//	return nil, err
	//}
	//
	//_, err = resourceContainer.RunGet(*version, params, volumeForGet)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return volumeForGet, nil
	return nil, nil
}

func NewCheckCommandRunner(r Resource, c garden.Container) (*checkCommandProcess, error) {
	p := NewCheckCommandProcess(c)

	i := checkRequest{
		Source:  r.Source,
		Version: nil,
	}

	input, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	var (
		stdout *bytes.Buffer
		stderr *bytes.Buffer
	)

	p.Stdin = bytes.NewBuffer(input)
	p.Stdout = stdout
	p.Stderr = stderr

	return &checkCommandProcess{
		ContainerProcess: p,
		out:              stdout,
		err:              stderr,
	}, nil
}

type checkCommandProcess struct {
	*ContainerProcess

	out *bytes.Buffer
	err *bytes.Buffer
}

func (c *checkCommandProcess) Versions() ([]atc.Version, error) {
	var o []atc.Version

	err := json.NewDecoder(c.out).Decode(&o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (c *checkCommandProcess) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	return c.ContainerProcess.Run(signals, ready)
}

type ContainerProcess struct {
	Container garden.Container

	Path   string
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func NewCheckCommandProcess(container garden.Container) *ContainerProcess {
	return &ContainerProcess{
		Container: container,
		Path:      "/opt/resource/check",
	}
}

func (c *ContainerProcess) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	processIO := garden.ProcessIO{
		Stdin:  c.Stdin,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
	}

	var process garden.Process

	process, err := c.Container.Run(garden.ProcessSpec{
		Path: c.Path,
		Args: c.Args,
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
			return container.ErrScriptFailed{
				Path:       c.Path,
				Args:       c.Args,
				ExitStatus: processStatus,
			}
		}

	case <-signals:
		c.Container.Stop(false)
		<-processExited
		return container.ErrAborted
	}

	return nil
}
