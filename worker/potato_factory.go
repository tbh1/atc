package worker

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
)

type PotatoFactory interface {
	NewBuildPotato(
		planID atc.PlanID,
		buildID int,
		pipelineID int,
	) Potato

	NewResourcePotato(
		resourceID int,
		pipelineID int,
	) Potato

	NewResourceTypePotato(
		resourceTypeID int,
		pipelineID int,
	) Potato
}

type Potato interface {
	ResourceInstance(
		imageResource atc.ImageResource,
		imageResourceVersion atc.Version,
		customTypes atc.ResourceTypes,
	) ResourceInstance
}

type ResourceInstance interface {
	FindOn(lager.Logger, Client) (Volume, bool, error)
	FindOrCreateOn(lager.Logger, Client) (Volume, error)

	ResourceCacheIdentifier() ResourceCacheIdentifier
}
