package resource_test

import (
	"github.com/concourse/atc"
	. "github.com/concourse/atc/cessna/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Check for new versions of resources", func() {

	var checkVersions []atc.Version
	var checkErr error

	Context("whose type is a base resource type", func() {
		BeforeEach(func() {
			source := atc.Source{
				"versions": []map[string]string{
					{"ref": "123"},
					{"beep": "boop"},
				},
			}

			testBaseResource = NewResource(baseResourceType, source)
		})

		JustBeforeEach(func() {
			checkVersions, checkErr = testBaseResource.Check(testWorker, nil)
		})

		It("runs the check script", func() {
			Expect(checkErr).ShouldNot(HaveOccurred())
		})

		It("returns the proper versions", func() {
			Expect(checkVersions).To(ConsistOf(atc.Version{"ref": "123"}, atc.Version{"beep": "boop"}))
		})

	})

})
