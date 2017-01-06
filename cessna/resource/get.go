package resource

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/cessna"
	"github.com/tedsuo/ifrit"
)

type getRequest struct {
	Source  atc.Source  `json:"source"`
	Params  atc.Params  `json:"params,omitempty"`
	Version atc.Version `json:"version,omitempty"`
}

type versionResult struct {
	Version  atc.Version         `json:"version"`
	Metadata []atc.MetadataField `json:"metadata,omitempty"`
}

func (rc *resourceContainer) RunGet(version atc.Version, params atc.Params, volume cessna.Volume) (versionResult, error) {
	var vr versionResult

	runner := rc.RunScript("/opt/resource/in", []string{"/tmp/resource/get"}, getRequest{rc.resource.Source, params, version}, &vr)

	getting := ifrit.Invoke(runner)

	err := <-getting.Wait()
	if err != nil {
		return versionResult{}, err
	}

	return vr, nil
}

// func (rc *resourceContainer) In() {
//
// }
//
// func (rc *resourceContainer) Out() {
//
// }
