package resource

import (
	"archive/tar"
	"bytes"
	"path"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/cessna"
	"github.com/concourse/baggageclaim"
	"github.com/tedsuo/ifrit"
)

type ResourceManager struct {
	Logger lager.Logger
	Worker *cessna.Worker
}

func NewResourceManagerFor(worker *cessna.Worker) *ResourceManager {
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
	rootFSVolume, err := r.Worker.BaggageClaimClient().CreateVolume(r.Logger.Session("create-cow-rootfs-volume"), parentVolume.Handle(), spec)
	if err != nil {
		return nil, err
	}

	// Turn RootFS COW into Container
	gardenSpec := garden.ContainerSpec{
		Privileged: false,
		RootFSPath: rootFSVolume.Path(),
	}

	gardenContainer, err := r.Worker.GardenClient().Create(gardenSpec)
	if err != nil {
		r.Logger.Error("failed-to-create-gardenContainer-in-garden", err)
		return nil, err
	}

	runner, err := NewCheckCommandProcess(gardenContainer, resource, version)
	if err != nil {
		return nil, err
	}

	checking := ifrit.Invoke(runner)

	err = <-checking.Wait()
	if err != nil {
		return nil, err
	}

	return runner.Response()
}

func (r *ResourceManager) Get(resource Resource, version *atc.Version, params atc.Params) (baggageclaim.Volume, error) {
	// Import RootFS into Volume
	spec := baggageclaim.VolumeSpec{
		Strategy: baggageclaim.ImportStrategy{
			Path: resource.ResourceType.RootFSPath,
		},
		Privileged: true,
	}
	handle := "foobar"

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
	rootFSVolume, err := r.Worker.BaggageClaimClient().CreateVolume(r.Logger.Session("create-cow-rootfs-volume"), parentVolume.Handle(), spec)
	if err != nil {
		return nil, err
	}

	// Empty Volume for Get
	spec = baggageclaim.VolumeSpec{
		Strategy:   baggageclaim.EmptyStrategy{},
		Privileged: false,
	}
	volumeForGet, err := r.Worker.BaggageClaimClient().CreateVolume(r.Logger.Session("create-empty-volume-for-get"), "foobar3", spec)
	if err != nil {
		return nil, err

	}

	// Turn into Container
	mount := garden.BindMount{
		SrcPath: volumeForGet.Path(),
		DstPath: "/tmp/resource/get",
		Mode:    garden.BindMountModeRW,
	}
	bindMounts := []garden.BindMount{mount}

	gardenSpec := garden.ContainerSpec{
		Privileged: false,
		RootFSPath: rootFSVolume.Path(),
		BindMounts: bindMounts,
	}

	container, err := r.Worker.GardenClient().Create(gardenSpec)
	if err != nil {
		r.Logger.Error("failed-to-create-gardenContainer-in-garden", err)
		return nil, err
	}

	runner, err := NewGetCommandProcess(container, mount, resource, version, params)
	if err != nil {
		r.Logger.Error("failed-to-create-get-command-process", err)
	}

	getting := ifrit.Invoke(runner)

	err = <-getting.Wait()
	if err != nil {
		return nil, err
	}

	return volumeForGet, nil
}

func (r *ResourceManager) Put(resource Resource, params atc.Params, artifacts NamedArtifacts) (OutResponse, error) {
	// Import RootFS into Volume
	spec := baggageclaim.VolumeSpec{
		Strategy: baggageclaim.ImportStrategy{
			Path: resource.ResourceType.RootFSPath,
		},
		Privileged: true,
	}
	handle := "foobar"

	parentVolume, err := r.Worker.BaggageClaimClient().CreateVolume(r.Logger.Session("create-parent-volume"), handle, spec)
	if err != nil {
		return OutResponse{}, err
	}

	// COW of RootFS Volume
	spec = baggageclaim.VolumeSpec{
		Strategy: baggageclaim.COWStrategy{
			Parent: parentVolume,
		},
		Privileged: false,
	}
	rootFSVolume, err := r.Worker.BaggageClaimClient().CreateVolume(r.Logger.Session("create-cow-rootfs-volume"), parentVolume.Handle(), spec)
	if err != nil {
		return OutResponse{}, err
	}

	// Turning artifacts into COWs
	cowArtifacts := make(NamedArtifacts)

	for name, volume := range artifacts {
		spec = baggageclaim.VolumeSpec{
			Strategy: baggageclaim.COWStrategy{
				Parent: volume,
			},
			Privileged: false,
		}
		v, err := r.Worker.BaggageClaimClient().CreateVolume(r.Logger.Session("create-cow-of-input-volume"), volume.Handle(), spec)
		if err != nil {
			return OutResponse{}, err
		}

		cowArtifacts[name] = v
	}

	// Create bindmounts for those COWs
	var bindMounts []garden.BindMount

	baseDirectory := "/tmp/artifacts"
	for name, volume := range cowArtifacts {
		bindMounts = append(bindMounts, garden.BindMount{
			SrcPath: volume.Path(),
			DstPath: path.Join(baseDirectory, name),
			Mode:    garden.BindMountModeRW,
		})
	}

	// Create container
	gardenSpec := garden.ContainerSpec{
		Privileged: false,
		RootFSPath: rootFSVolume.Path(),
		BindMounts: bindMounts,
	}

	container, err := r.Worker.GardenClient().Create(gardenSpec)
	if err != nil {
		r.Logger.Error("failed-to-create-gardenContainer-in-garden", err)
		return OutResponse{}, err
	}

	// Stream fake tar into container to make sure directory exists
	// stream into baseDirectory
	emptyTar := new(bytes.Buffer)

	err = tar.NewWriter(emptyTar).Close()
	if err != nil {
		return OutResponse{}, err
	}

	err = container.StreamIn(garden.StreamInSpec{
		Path:      baseDirectory,
		TarStream: emptyTar,
	})

	// Create the PutProcess
	runner, err := NewPutCommandProcess(container, resource, params, baseDirectory)
	if err != nil {
		r.Logger.Error("failed-to-create-get-command-process", err)
	}

	// Run the PutProcess
	putting := ifrit.Invoke(runner)

	err = <-putting.Wait()
	if err != nil {
		return OutResponse{}, err
	}

	// Parse the PutProcess output
	return runner.Response()
}
