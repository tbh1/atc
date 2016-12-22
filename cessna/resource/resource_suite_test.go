package resource_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/concourse/atc/cessna"
	. "github.com/concourse/atc/cessna/resource"

	"testing"
)

var testBaseResource Resource
var testWorker cessna.Worker
var baseResourceType ResourceType

func TestResource(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		testWorker = cessna.Worker{
			GardenAddr:       "10.244.16.2:7777",
			BaggageclaimAddr: "10.244.16.2:7788",
		}

		baseResourceType = ResourceType{
			RootFSPath: "/var/vcap/data/packages/echo_resource/rootfs",
			Name:       "git",
		}

	})

	RunSpecs(t, "Resource Suite")

}
