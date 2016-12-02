package potato

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/worker"
)

type potatoFactory struct {
	resourceInstanceFactory resource.ResourceInstanceFactory
}

func NewPotatoFactory(resourceInstanceFactory resource.ResourceInstanceFactory) worker.PotatoFactory {
	return &potatoFactory{
		resourceInstanceFactory: resourceInstanceFactory,
	}
}

type buildPotato struct {
	planID     atc.PlanID
	buildID    int
	pipelineID int

	resourceInstanceFactory resource.ResourceInstanceFactory
}

func (f *potatoFactory) NewBuildPotato(
	planID atc.PlanID,
	buildID int,
	pipelineID int,
) worker.Potato {
	return &buildPotato{
		planID:                  planID,
		buildID:                 buildID,
		pipelineID:              pipelineID,
		resourceInstanceFactory: f.resourceInstanceFactory,
	}
}

func (p *buildPotato) ResourceInstance(
	imageResource atc.ImageResource,
	imageResourceVersion atc.Version,
	customTypes atc.ResourceTypes,
) worker.ResourceInstance {
	return p.resourceInstanceFactory.NewBuildResourceInstance(
		resource.ResourceType(imageResource.Type),
		imageResourceVersion,
		imageResource.Source,
		nil,
		&dbng.Build{ID: p.buildID},
		&dbng.Pipeline{ID: p.pipelineID},
		customTypes,
	)
}

type resourcePotato struct {
	resourceID int
	pipelineID int

	resourceInstanceFactory resource.ResourceInstanceFactory
}

func (f *potatoFactory) NewResourcePotato(
	resourceID int,
	pipelineID int,
) worker.Potato {
	return &resourcePotato{
		resourceID:              resourceID,
		pipelineID:              pipelineID,
		resourceInstanceFactory: f.resourceInstanceFactory,
	}
}

func (p *resourcePotato) ResourceInstance(
	imageResource atc.ImageResource,
	imageResourceVersion atc.Version,
	customTypes atc.ResourceTypes,
) worker.ResourceInstance {
	return p.resourceInstanceFactory.NewResourceResourceInstance(
		resource.ResourceType(imageResource.Type),
		imageResourceVersion,
		imageResource.Source,
		nil,
		&dbng.Resource{ID: p.resourceID},
		&dbng.Pipeline{ID: p.pipelineID},
		customTypes,
	)
}

type resourceTypePotato struct {
	resourceTypeID int
	pipelineID     int

	resourceInstanceFactory resource.ResourceInstanceFactory
}

func (f *potatoFactory) NewResourceTypePotato(
	resourceTypeID int,
	pipelineID int,
) worker.Potato {
	return &resourceTypePotato{
		resourceTypeID:          resourceTypeID,
		pipelineID:              pipelineID,
		resourceInstanceFactory: f.resourceInstanceFactory,
	}
}

func (p *resourceTypePotato) ResourceInstance(
	imageResource atc.ImageResource,
	imageResourceVersion atc.Version,
	customTypes atc.ResourceTypes,
) worker.ResourceInstance {
	return p.resourceInstanceFactory.NewResourceTypeResourceInstance(
		resource.ResourceType(imageResource.Type),
		imageResourceVersion,
		imageResource.Source,
		nil,
		&dbng.UsedResourceType{ID: p.resourceTypeID},
		&dbng.Pipeline{ID: p.pipelineID},
		customTypes,
	)
}
