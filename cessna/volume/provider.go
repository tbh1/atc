package volume

import (
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/cessna"

	"github.com/concourse/baggageclaim"
)

type Provider struct {
	baggageclaimClient baggageclaim.Client
	logger             lager.Logger
}

func NewProvider(worker cessna.Worker) Provider {
	return Provider{
		baggageclaimClient: worker.BaggageClaimClient(),
		logger:             lager.NewLogger("worker-volume-provider"),
	}
}

func (p *Provider) CreateEmptyVolume(handle string, privileged bool) (cessna.Volume, error) {
	spec := baggageclaim.VolumeSpec{
		Strategy:   baggageclaim.EmptyStrategy{},
		Privileged: privileged,
	}
	return p.createVolume(handle, spec)
}

func (p *Provider) ImportRootFS(path string, privileged bool) (cessna.Volume, error) {
	spec := baggageclaim.VolumeSpec{
		Strategy: baggageclaim.ImportStrategy{
			Path: path,
		},
		Privileged: privileged,
	}
	return p.createVolume(path, spec)
}

func (p *Provider) COWFromRootFS(rootFSPath string) (cessna.Volume, error) {
	vol, err := p.ImportRootFS(rootFSPath, true)
	if err != nil {
		p.logger.Error("import-rootfs-failed", err)
		return nil, err
	}

	cowVolume, err := vol.COWify(false)
	if err != nil {
		p.logger.Error("cow-volume-creation-failed", err)
		return nil, err
	}

	return cowVolume, err
}

func (p *Provider) createVolume(handle string, spec baggageclaim.VolumeSpec) (cessna.Volume, error) {
	baggageclaimVolume, err := p.baggageclaimClient.CreateVolume(p.logger.Session("baggageclaim-create-volume"), handle, spec)
	return &workerVolume{baggageclaimVolume, p}, err
}

///an existing volume on a worker (one that can be COW'd by baggageclaim)
type workerVolume struct {
	baggageclaim.Volume
	provider *Provider
}

func (wv *workerVolume) COWify(privileged bool) (cessna.Volume, error) {
	spec := baggageclaim.VolumeSpec{
		Strategy:   baggageclaim.COWStrategy{Parent: wv.Volume},
		Privileged: privileged,
	}

	return wv.provider.createVolume(fmt.Sprintf("cow-%s", wv.Handle()), spec)
}
