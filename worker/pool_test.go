package worker_test

import (
	"errors"
	"os"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/dbng"
	. "github.com/concourse/atc/worker"
	"github.com/concourse/atc/worker/workerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pool", func() {
	var (
		logger       *lagertest.TestLogger
		fakeProvider *workerfakes.FakeWorkerProvider

		pool Client
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		fakeProvider = new(workerfakes.FakeWorkerProvider)

		pool = NewPool(fakeProvider)
	})

	Describe("GetWorker", func() {
		Context("when the call to lookup the worker returns an error", func() {
			BeforeEach(func() {
				fakeProvider.GetWorkerReturns(nil, false, errors.New("disaster"))
			})

			It("returns an error", func() {
				foundWorker, err := pool.GetWorker("some-worker")
				Expect(err).To(HaveOccurred())
				Expect(foundWorker).To(BeNil())
			})
		})

		Context("when the call to lookup the worker fails because the worker was not found", func() {
			BeforeEach(func() {
				fakeProvider.GetWorkerReturns(nil, false, nil)
			})

			It("returns an error indicating no workers were found", func() {
				foundWorker, err := pool.GetWorker("no-worker")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrNoWorkers))
				Expect(foundWorker).To(BeNil())
			})
		})

		Context("when the lookup of the worker succeeds", func() {
			var fakeWorker *workerfakes.FakeWorker

			BeforeEach(func() {
				fakeWorker = new(workerfakes.FakeWorker)
				fakeProvider.GetWorkerReturns(fakeWorker, true, nil)
			})

			It("returns an error indicating no workers were found", func() {
				foundWorker, err := pool.GetWorker("some-worker")
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeProvider.GetWorkerCallCount()).To(Equal(1))
				workerName := fakeProvider.GetWorkerArgsForCall(0)
				Expect(workerName).To(Equal("some-worker"))
				Expect(foundWorker).To(Equal(fakeWorker))
			})
		})
	})

	Describe("Satisfying", func() {
		var (
			spec WorkerSpec

			satisfyingErr    error
			satisfyingWorker Worker
			resourceTypes    atc.VersionedResourceTypes
		)

		BeforeEach(func() {
			spec = WorkerSpec{
				Platform: "some-platform",
				Tags:     []string{"step", "tags"},
			}
			resourceTypes = atc.VersionedResourceTypes{
				{
					ResourceType: atc.ResourceType{
						Name:   "some-resource-type",
						Type:   "some-underlying-type",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
			}
		})

		JustBeforeEach(func() {
			satisfyingWorker, satisfyingErr = pool.Satisfying(spec, resourceTypes)
		})

		Context("with multiple workers", func() {
			var (
				workerA *workerfakes.FakeWorker
				workerB *workerfakes.FakeWorker
				workerC *workerfakes.FakeWorker
			)

			BeforeEach(func() {
				workerA = new(workerfakes.FakeWorker)
				workerB = new(workerfakes.FakeWorker)
				workerC = new(workerfakes.FakeWorker)

				workerA.SatisfyingReturns(workerA, nil)
				workerB.SatisfyingReturns(workerB, nil)
				workerC.SatisfyingReturns(nil, errors.New("nope"))

				fakeProvider.RunningWorkersReturns([]Worker{workerA, workerB, workerC}, nil)
			})

			It("succeeds", func() {
				Expect(satisfyingErr).NotTo(HaveOccurred())
			})

			It("checks that the workers satisfy the given spec", func() {
				Expect(workerA.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes := workerA.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec))
				Expect(actualResourceTypes).To(Equal(resourceTypes))

				Expect(workerB.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes = workerB.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec))
				Expect(actualResourceTypes).To(Equal(resourceTypes))

				Expect(workerC.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes = workerC.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec))
				Expect(actualResourceTypes).To(Equal(resourceTypes))
			})

			It("returns a random worker satisfying the spec", func() {
				chosenCount := map[Worker]int{workerA: 0, workerB: 0, workerC: 0}
				for i := 0; i < 100; i++ {
					satisfyingWorker, satisfyingErr = pool.Satisfying(spec, resourceTypes)
					Expect(satisfyingErr).NotTo(HaveOccurred())
					chosenCount[satisfyingWorker]++
				}
				Expect(chosenCount[workerA]).To(BeNumerically("~", chosenCount[workerB], 50))
				Expect(chosenCount[workerC]).To(BeZero())
			})

			Context("when no workers satisfy the spec", func() {
				BeforeEach(func() {
					workerA.SatisfyingReturns(nil, errors.New("nope"))
					workerB.SatisfyingReturns(nil, errors.New("nope"))
					workerC.SatisfyingReturns(nil, errors.New("nope"))
				})

				It("returns a NoCompatibleWorkersError", func() {
					Expect(satisfyingErr).To(Equal(NoCompatibleWorkersError{
						Spec:    spec,
						Workers: []Worker{workerA, workerB, workerC},
					}))
				})
			})
		})

		Context("with no workers", func() {
			BeforeEach(func() {
				fakeProvider.RunningWorkersReturns([]Worker{}, nil)
			})

			It("returns ErrNoWorkers", func() {
				Expect(satisfyingErr).To(Equal(ErrNoWorkers))
			})
		})

		Context("when getting the workers fails", func() {
			disaster := errors.New("nope")

			BeforeEach(func() {
				fakeProvider.RunningWorkersReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(satisfyingErr).To(Equal(disaster))
			})
		})
	})

	Describe("AllSatisfying", func() {
		var (
			spec WorkerSpec

			satisfyingErr     error
			satisfyingWorkers []Worker
			resourceTypes     atc.VersionedResourceTypes
		)

		BeforeEach(func() {
			spec = WorkerSpec{
				Platform: "some-platform",
				Tags:     []string{"step", "tags"},
			}
			resourceTypes = atc.VersionedResourceTypes{
				{
					ResourceType: atc.ResourceType{
						Name:   "some-resource-type",
						Type:   "some-underlying-type",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
			}
		})

		JustBeforeEach(func() {
			satisfyingWorkers, satisfyingErr = pool.AllSatisfying(spec, resourceTypes)
		})

		Context("with multiple workers", func() {
			var (
				workerA *workerfakes.FakeWorker
				workerB *workerfakes.FakeWorker
				workerC *workerfakes.FakeWorker
			)

			BeforeEach(func() {
				workerA = new(workerfakes.FakeWorker)
				workerB = new(workerfakes.FakeWorker)
				workerC = new(workerfakes.FakeWorker)

				workerA.SatisfyingReturns(workerA, nil)
				workerB.SatisfyingReturns(workerB, nil)
				workerC.SatisfyingReturns(nil, errors.New("nope"))

				fakeProvider.RunningWorkersReturns([]Worker{workerA, workerB, workerC}, nil)
			})

			It("succeeds", func() {
				Expect(satisfyingErr).NotTo(HaveOccurred())
			})

			It("checks that the workers satisfy the given spec", func() {
				Expect(workerA.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes := workerA.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec))
				Expect(actualResourceTypes).To(Equal(resourceTypes))

				Expect(workerB.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes = workerB.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec))
				Expect(actualResourceTypes).To(Equal(resourceTypes))

				Expect(workerC.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes = workerC.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec))
				Expect(actualResourceTypes).To(Equal(resourceTypes))
			})

			It("returns all workers satisfying the spec in a random order", func() {
				firstCount := map[Worker]int{workerA: 0, workerB: 0}
				for i := 0; i < 100; i++ {
					satisfyingWorkers, satisfyingErr = pool.AllSatisfying(spec, resourceTypes)
					Expect(satisfyingErr).NotTo(HaveOccurred())
					Expect(satisfyingWorkers).To(ConsistOf(workerA, workerB))
					firstCount[satisfyingWorkers[0]]++
				}
				Expect(firstCount[workerA]).To(BeNumerically("~", firstCount[workerB], 50))
			})

			Context("when no workers satisfy the spec", func() {
				BeforeEach(func() {
					workerA.SatisfyingReturns(nil, errors.New("nope"))
					workerB.SatisfyingReturns(nil, errors.New("nope"))
					workerC.SatisfyingReturns(nil, errors.New("nope"))
				})

				It("returns a NoCompatibleWorkersError", func() {
					Expect(satisfyingErr).To(Equal(NoCompatibleWorkersError{
						Spec:    spec,
						Workers: []Worker{workerA, workerB, workerC},
					}))
				})
			})
		})

		Context("when team workers and general workers satisfy the spec", func() {
			var (
				teamWorker1   *workerfakes.FakeWorker
				teamWorker2   *workerfakes.FakeWorker
				teamWorker3   *workerfakes.FakeWorker
				generalWorker *workerfakes.FakeWorker
			)

			BeforeEach(func() {
				teamWorker1 = new(workerfakes.FakeWorker)
				teamWorker1.SatisfyingReturns(teamWorker1, nil)
				teamWorker1.IsOwnedByTeamReturns(true)
				teamWorker2 = new(workerfakes.FakeWorker)
				teamWorker2.SatisfyingReturns(teamWorker2, nil)
				teamWorker2.IsOwnedByTeamReturns(true)
				teamWorker3 = new(workerfakes.FakeWorker)
				teamWorker3.SatisfyingReturns(nil, errors.New("nope"))
				generalWorker = new(workerfakes.FakeWorker)
				generalWorker.SatisfyingReturns(generalWorker, nil)
				generalWorker.IsOwnedByTeamReturns(false)
				fakeProvider.RunningWorkersReturns([]Worker{generalWorker, teamWorker1, teamWorker2, teamWorker3}, nil)
			})

			It("returns only the team workers that satisfy the spec", func() {
				Expect(satisfyingErr).NotTo(HaveOccurred())
				Expect(satisfyingWorkers).To(ConsistOf(teamWorker1, teamWorker2))
			})
		})

		Context("when only general workers satisfy the spec", func() {
			var (
				teamWorker     *workerfakes.FakeWorker
				generalWorker1 *workerfakes.FakeWorker
				generalWorker2 *workerfakes.FakeWorker
			)

			BeforeEach(func() {
				teamWorker = new(workerfakes.FakeWorker)
				teamWorker.SatisfyingReturns(nil, errors.New("nope"))
				generalWorker1 = new(workerfakes.FakeWorker)
				generalWorker1.SatisfyingReturns(generalWorker1, nil)
				generalWorker1.IsOwnedByTeamReturns(false)
				generalWorker2 = new(workerfakes.FakeWorker)
				generalWorker2.SatisfyingReturns(nil, errors.New("nope"))
				fakeProvider.RunningWorkersReturns([]Worker{generalWorker1, generalWorker2, teamWorker}, nil)
			})

			It("returns the general workers that satisfy the spec", func() {
				Expect(satisfyingErr).NotTo(HaveOccurred())
				Expect(satisfyingWorkers).To(ConsistOf(generalWorker1))
			})
		})

		Context("with no workers", func() {
			BeforeEach(func() {
				fakeProvider.RunningWorkersReturns([]Worker{}, nil)
			})

			It("returns ErrNoWorkers", func() {
				Expect(satisfyingErr).To(Equal(ErrNoWorkers))
			})
		})

		Context("when getting the workers fails", func() {
			disaster := errors.New("nope")

			BeforeEach(func() {
				fakeProvider.RunningWorkersReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(satisfyingErr).To(Equal(disaster))
			})
		})
	})

	Describe("FindOrCreateBuildContainer", func() {
		var (
			fakeImageFetchingDelegate *workerfakes.FakeImageFetchingDelegate

			id   Identifier
			spec ContainerSpec

			fakeContainer *workerfakes.FakeContainer

			createdContainer Container
			createErr        error
			resourceTypes    atc.VersionedResourceTypes
		)

		BeforeEach(func() {
			fakeImageFetchingDelegate = new(workerfakes.FakeImageFetchingDelegate)

			id = Identifier{
				ResourceID: 1234,
			}

			spec = ContainerSpec{
				ImageSpec: ImageSpec{ResourceType: "some-type"},
				TeamID:    4567,
			}

			resourceTypes = atc.VersionedResourceTypes{
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-b",
						Type:   "custom-type-a",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-a",
						Type:   "some-resource",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-c",
						Type:   "custom-type-b",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-d",
						Type:   "custom-type-b",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "unknown-custom-type",
						Type:   "unknown-base-type",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
			}

			fakeContainer = new(workerfakes.FakeContainer)
		})

		JustBeforeEach(func() {
			createdContainer, createErr = pool.FindOrCreateBuildContainer(
				logger,
				make(chan os.Signal),
				fakeImageFetchingDelegate,
				id,
				Metadata{},
				spec,
				resourceTypes,
				map[string]string{},
			)
		})

		Context("when a worker is found with the container", func() {
			var fakeWorker *workerfakes.FakeWorker

			BeforeEach(func() {
				fakeWorker = new(workerfakes.FakeWorker)
				fakeProvider.FindWorkerForBuildContainerReturns(fakeWorker, true, nil)
				fakeWorker.FindOrCreateBuildContainerReturns(fakeContainer, nil)
			})

			It("succeeds", func() {
				Expect(createErr).NotTo(HaveOccurred())
			})

			It("returns the created container", func() {
				Expect(createdContainer).To(Equal(fakeContainer))
			})

			It("'find-or-create's on the particular worker", func() {
				_, actualTeamID, actualResourceUser, actualResourceType, actualResourceSource, actualResourceTypes := fakeProvider.FindWorkerForBuildContainerArgsForCall(0)
				Expect(actualTeamID).To(Equal(4567))
				Expect(actualResourceUser).To(Equal(dbng.ForBuild{42}))
				Expect(actualResourceType).To(Equal("some-type"))
				Expect(actualResourceSource).To(Equal(atc.Source{"some": "source"}))
				Expect(actualResourceTypes).To(Equal(resourceTypes))

				Expect(fakeWorker.FindOrCreateBuildContainerCallCount()).To(Equal(1))
			})
		})

		Context("when a worker is not found, and multiple are present", func() {
			var (
				workerA *workerfakes.FakeWorker
				workerB *workerfakes.FakeWorker
				workerC *workerfakes.FakeWorker
			)

			BeforeEach(func() {
				fakeProvider.FindWorkerForBuildContainerReturns(nil, false, nil)

				workerA = new(workerfakes.FakeWorker)
				workerB = new(workerfakes.FakeWorker)
				workerC = new(workerfakes.FakeWorker)

				workerA.ActiveContainersReturns(3)
				workerB.ActiveContainersReturns(2)

				workerA.SatisfyingReturns(workerA, nil)
				workerB.SatisfyingReturns(workerB, nil)
				workerC.SatisfyingReturns(nil, errors.New("nope"))

				workerA.FindOrCreateBuildContainerReturns(fakeContainer, nil)
				workerB.FindOrCreateBuildContainerReturns(fakeContainer, nil)
				workerC.FindOrCreateBuildContainerReturns(fakeContainer, nil)

				fakeProvider.RunningWorkersReturns([]Worker{workerA, workerB, workerC}, nil)
			})

			It("succeeds", func() {
				Expect(createErr).NotTo(HaveOccurred())
			})

			It("returns the created container", func() {
				Expect(createdContainer).To(Equal(fakeContainer))
			})

			It("checks that the workers satisfy the given spec", func() {
				Expect(workerA.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes := workerA.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec.WorkerSpec()))
				Expect(actualResourceTypes).To(Equal(resourceTypes))

				Expect(workerB.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes = workerB.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec.WorkerSpec()))
				Expect(actualResourceTypes).To(Equal(resourceTypes))

				Expect(workerC.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes = workerC.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec.WorkerSpec()))
				Expect(actualResourceTypes).To(Equal(resourceTypes))
			})

			It("creates using a random worker", func() {
				for i := 1; i < 100; i++ { // account for initial create in JustBefore
					createdContainer, createErr := pool.FindOrCreateBuildContainer(
						logger,
						make(chan os.Signal),
						fakeImageFetchingDelegate,
						id,
						Metadata{},
						spec,
						resourceTypes,
						map[string]string{},
					)
					Expect(createErr).NotTo(HaveOccurred())
					Expect(createdContainer).To(Equal(fakeContainer))
				}

				Expect(workerA.FindOrCreateBuildContainerCallCount()).To(BeNumerically("~", workerB.FindOrCreateBuildContainerCallCount(), 50))
				Expect(workerC.FindOrCreateBuildContainerCallCount()).To(BeZero())
			})

			Context("when creating the container fails", func() {
				disaster := errors.New("nope")

				BeforeEach(func() {
					workerA.FindOrCreateBuildContainerReturns(nil, disaster)
					workerB.FindOrCreateBuildContainerReturns(nil, disaster)
				})

				It("returns the error", func() {
					Expect(createErr).To(Equal(disaster))
				})
			})

			Context("when no workers satisfy the spec", func() {
				BeforeEach(func() {
					workerA.SatisfyingReturns(nil, errors.New("nope"))
					workerB.SatisfyingReturns(nil, errors.New("nope"))
					workerC.SatisfyingReturns(nil, errors.New("nope"))
				})

				It("returns a NoCompatibleWorkersError", func() {
					Expect(createErr).To(Equal(NoCompatibleWorkersError{
						Spec:    spec.WorkerSpec(),
						Workers: []Worker{workerA, workerB, workerC},
					}))

				})
			})
		})

		Context("with no workers", func() {
			BeforeEach(func() {
				fakeProvider.RunningWorkersReturns([]Worker{}, nil)
			})

			It("returns ErrNoWorkers", func() {
				Expect(createErr).To(Equal(ErrNoWorkers))
			})
		})

		Context("when getting the workers fails", func() {
			disaster := errors.New("nope")

			BeforeEach(func() {
				fakeProvider.RunningWorkersReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(createErr).To(Equal(disaster))
			})
		})
	})

	Describe("FindOrCreateResourceCheckContainer", func() {
		var (
			fakeImageFetchingDelegate *workerfakes.FakeImageFetchingDelegate

			id   Identifier
			spec ContainerSpec

			fakeContainer *workerfakes.FakeContainer

			createdContainer Container
			createErr        error
			resourceTypes    atc.VersionedResourceTypes
		)

		BeforeEach(func() {
			fakeImageFetchingDelegate = new(workerfakes.FakeImageFetchingDelegate)

			id = Identifier{
				ResourceID: 1234,
			}

			spec = ContainerSpec{
				ImageSpec: ImageSpec{ResourceType: "some-type"},
				TeamID:    4567,
			}

			resourceTypes = atc.VersionedResourceTypes{
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-b",
						Type:   "custom-type-a",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-a",
						Type:   "some-resource",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-c",
						Type:   "custom-type-b",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-d",
						Type:   "custom-type-b",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "unknown-custom-type",
						Type:   "unknown-base-type",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
			}

			fakeContainer = new(workerfakes.FakeContainer)
		})

		JustBeforeEach(func() {
			createdContainer, createErr = pool.FindOrCreateResourceCheckContainer(
				logger,
				dbng.ForBuild{42},
				make(chan os.Signal),
				fakeImageFetchingDelegate,
				id,
				Metadata{},
				spec,
				resourceTypes,
				"some-type",
				atc.Source{"some": "source"},
			)
		})

		Context("when a worker is found with the container", func() {
			var fakeWorker *workerfakes.FakeWorker

			BeforeEach(func() {
				fakeWorker = new(workerfakes.FakeWorker)
				fakeProvider.FindWorkerForResourceCheckContainerReturns(fakeWorker, true, nil)
				fakeWorker.FindOrCreateResourceCheckContainerReturns(fakeContainer, nil)
			})

			It("succeeds", func() {
				Expect(createErr).NotTo(HaveOccurred())
			})

			It("returns the created container", func() {
				Expect(createdContainer).To(Equal(fakeContainer))
			})

			It("'find-or-create's on the particular worker", func() {
				_, actualTeamID, actualResourceUser, actualResourceType, actualResourceSource, actualResourceTypes := fakeProvider.FindWorkerForResourceCheckContainerArgsForCall(0)
				Expect(actualTeamID).To(Equal(4567))
				Expect(actualResourceUser).To(Equal(dbng.ForBuild{42}))
				Expect(actualResourceType).To(Equal("some-type"))
				Expect(actualResourceSource).To(Equal(atc.Source{"some": "source"}))
				Expect(actualResourceTypes).To(Equal(resourceTypes))

				Expect(fakeWorker.FindOrCreateResourceCheckContainerCallCount()).To(Equal(1))
			})
		})

		Context("when a worker is not found, and multiple are present", func() {
			var (
				workerA *workerfakes.FakeWorker
				workerB *workerfakes.FakeWorker
				workerC *workerfakes.FakeWorker
			)

			BeforeEach(func() {
				fakeProvider.FindWorkerForResourceCheckContainerReturns(nil, false, nil)

				workerA = new(workerfakes.FakeWorker)
				workerB = new(workerfakes.FakeWorker)
				workerC = new(workerfakes.FakeWorker)

				workerA.ActiveContainersReturns(3)
				workerB.ActiveContainersReturns(2)

				workerA.SatisfyingReturns(workerA, nil)
				workerB.SatisfyingReturns(workerB, nil)
				workerC.SatisfyingReturns(nil, errors.New("nope"))

				workerA.FindOrCreateResourceCheckContainerReturns(fakeContainer, nil)
				workerB.FindOrCreateResourceCheckContainerReturns(fakeContainer, nil)
				workerC.FindOrCreateResourceCheckContainerReturns(fakeContainer, nil)

				fakeProvider.RunningWorkersReturns([]Worker{workerA, workerB, workerC}, nil)
			})

			It("succeeds", func() {
				Expect(createErr).NotTo(HaveOccurred())
			})

			It("returns the created container", func() {
				Expect(createdContainer).To(Equal(fakeContainer))
			})

			It("checks that the workers satisfy the given spec", func() {
				Expect(workerA.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes := workerA.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec.WorkerSpec()))
				Expect(actualResourceTypes).To(Equal(resourceTypes))

				Expect(workerB.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes = workerB.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec.WorkerSpec()))
				Expect(actualResourceTypes).To(Equal(resourceTypes))

				Expect(workerC.SatisfyingCallCount()).To(Equal(1))
				actualSpec, actualResourceTypes = workerC.SatisfyingArgsForCall(0)
				Expect(actualSpec).To(Equal(spec.WorkerSpec()))
				Expect(actualResourceTypes).To(Equal(resourceTypes))
			})

			It("creates using a random worker", func() {
				for i := 1; i < 100; i++ { // account for initial create in JustBefore
					createdContainer, createErr := pool.FindOrCreateResourceCheckContainer(
						logger,
						dbng.ForBuild{42},
						make(chan os.Signal),
						fakeImageFetchingDelegate,
						id,
						Metadata{},
						spec,
						resourceTypes,
						"some-type",
						atc.Source{"some": "source"},
					)
					Expect(createErr).NotTo(HaveOccurred())
					Expect(createdContainer).To(Equal(fakeContainer))
				}

				Expect(workerA.FindOrCreateResourceCheckContainerCallCount()).To(BeNumerically("~", workerB.FindOrCreateResourceCheckContainerCallCount(), 50))
				Expect(workerC.FindOrCreateResourceCheckContainerCallCount()).To(BeZero())
			})

			Context("when creating the container fails", func() {
				disaster := errors.New("nope")

				BeforeEach(func() {
					workerA.FindOrCreateResourceCheckContainerReturns(nil, disaster)
					workerB.FindOrCreateResourceCheckContainerReturns(nil, disaster)
				})

				It("returns the error", func() {
					Expect(createErr).To(Equal(disaster))
				})
			})

			Context("when no workers satisfy the spec", func() {
				BeforeEach(func() {
					workerA.SatisfyingReturns(nil, errors.New("nope"))
					workerB.SatisfyingReturns(nil, errors.New("nope"))
					workerC.SatisfyingReturns(nil, errors.New("nope"))
				})

				It("returns a NoCompatibleWorkersError", func() {
					Expect(createErr).To(Equal(NoCompatibleWorkersError{
						Spec:    spec.WorkerSpec(),
						Workers: []Worker{workerA, workerB, workerC},
					}))

				})
			})
		})

		Context("with no workers", func() {
			BeforeEach(func() {
				fakeProvider.RunningWorkersReturns([]Worker{}, nil)
			})

			It("returns ErrNoWorkers", func() {
				Expect(createErr).To(Equal(ErrNoWorkers))
			})
		})

		Context("when getting the workers fails", func() {
			disaster := errors.New("nope")

			BeforeEach(func() {
				fakeProvider.RunningWorkersReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(createErr).To(Equal(disaster))
			})
		})
	})

	Describe("FindContainerByHandle", func() {
		Context("when looking up the container info contains an error", func() {
			BeforeEach(func() {
				fakeProvider.GetContainerReturns(db.SavedContainer{}, false, errors.New("disaster"))
			})

			It("returns the error", func() {
				container, found, err := pool.FindContainerByHandle(logger, "some-handle", 42)
				Expect(err).To(HaveOccurred())
				Expect(container).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when looking up the container info does not find the container info", func() {
			BeforeEach(func() {
				fakeProvider.GetContainerReturns(db.SavedContainer{}, false, nil)
			})

			It("returns that it was not found", func() {
				container, found, err := pool.FindContainerByHandle(logger, "some-handle", 42)
				Expect(err).NotTo(HaveOccurred())
				Expect(container).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when looking up the container info is successful", func() {
			var container db.SavedContainer
			BeforeEach(func() {
				container = db.SavedContainer{
					Container: db.Container{
						ContainerMetadata: db.ContainerMetadata{
							WorkerName: "some-worker",
							Handle:     "some-container-handle",
						},
					},
				}

				fakeProvider.GetContainerReturns(container, true, nil)
			})

			It("calls to lookup the worker by name", func() {
				pool.FindContainerByHandle(logger, "some-container-handle", 42)

				Expect(fakeProvider.GetWorkerCallCount()).To(Equal(1))

				workerName := fakeProvider.GetWorkerArgsForCall(0)
				Expect(workerName).To(Equal("some-worker"))
			})

			Context("when looking up the worker returns an error", func() {
				BeforeEach(func() {
					fakeProvider.GetWorkerReturns(nil, false, errors.New("disaster"))
				})

				It("returns the error", func() {
					container, found, err := pool.FindContainerByHandle(logger, "some-handle", 42)
					Expect(err).To(HaveOccurred())
					Expect(container).To(BeNil())
					Expect(found).To(BeFalse())
				})
			})

			Context("when we cannot find the worker from the container info", func() {
				BeforeEach(func() {
					fakeProvider.GetWorkerReturns(nil, false, nil)
				})

				It("returns ErrMissingWorker", func() {
					container, found, err := pool.FindContainerByHandle(logger, "some-handle", 42)
					Expect(err).To(Equal(ErrMissingWorker))
					Expect(container).To(BeNil())
					Expect(found).To(BeFalse())
				})
			})

			Context("when looking up the worker is successful", func() {
				var fakeWorker *workerfakes.FakeWorker

				BeforeEach(func() {
					fakeWorker = new(workerfakes.FakeWorker)
					fakeProvider.GetWorkerReturns(fakeWorker, true, nil)
				})

				It("calls to lookup the container on the worker", func() {
					pool.FindContainerByHandle(logger, "some-handle", 42)

					Expect(fakeWorker.FindContainerByHandleCallCount()).To(Equal(1))

					_, handleArg, _ := fakeWorker.FindContainerByHandleArgsForCall(0)
					Expect(handleArg).To(Equal("some-handle"))
				})

				Context("when looking up the container contains an error", func() {
					It("returns the error", func() {
						fakeWorker.FindContainerByHandleReturns(nil, false, errors.New("disaster"))

						container, found, err := pool.FindContainerByHandle(logger, "some-handle", 42)
						Expect(err).To(HaveOccurred())
						Expect(container).To(BeNil())
						Expect(found).To(BeFalse())
					})
				})

				Context("when the finding the container on the worker is successful", func() {
					It("returns the container", func() {
						var fakeContainer *workerfakes.FakeContainer
						fakeContainer = new(workerfakes.FakeContainer)

						fakeWorker.FindContainerByHandleReturns(fakeContainer, true, nil)

						foundContainer, found, err := pool.FindContainerByHandle(logger, "some-handle", 42)
						Expect(err).NotTo(HaveOccurred())
						Expect(found).To(BeTrue())
						Expect(foundContainer).To(Equal(fakeContainer))
					})
				})
			})
		})
	})

	Describe("FindContainerForIdentifier", func() {
		var identifier Identifier

		BeforeEach(func() {
			identifier = Identifier{
				ResourceID: 1234,
			}
		})

		Context("when looking up the container info contains an error", func() {
			BeforeEach(func() {
				fakeProvider.FindContainerForIdentifierReturns(db.SavedContainer{}, false, errors.New("disaster"))
			})

			It("returns the error", func() {
				container, found, err := pool.FindContainerForIdentifier(logger, identifier)
				Expect(err).To(HaveOccurred())
				Expect(container).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when looking up the container info does not find the container info", func() {
			BeforeEach(func() {
				fakeProvider.FindContainerForIdentifierReturns(db.SavedContainer{}, false, nil)
			})

			It("returns that it was not found", func() {
				container, found, err := pool.FindContainerForIdentifier(logger, identifier)
				Expect(err).NotTo(HaveOccurred())
				Expect(container).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when looking up the container info is successful", func() {
			var container db.SavedContainer
			BeforeEach(func() {
				container = db.SavedContainer{
					Container: db.Container{
						ContainerMetadata: db.ContainerMetadata{
							WorkerName: "some-worker",
							Handle:     "some-container-handle",
							Type:       "checked",
						},
					},
				}

				fakeProvider.FindContainerForIdentifierReturns(container, true, nil)
			})

			It("calls to lookup the worker by name", func() {
				pool.FindContainerForIdentifier(logger, identifier)

				Expect(fakeProvider.GetWorkerCallCount()).To(Equal(1))

				workerName := fakeProvider.GetWorkerArgsForCall(0)
				Expect(workerName).To(Equal("some-worker"))
			})

			Context("when looking up the worker returns an error", func() {
				It("returns the error", func() {
					fakeProvider.GetWorkerReturns(nil, false, errors.New("disaster"))

					workerContainer, found, err := pool.FindContainerForIdentifier(logger, identifier)
					Expect(err).To(HaveOccurred())
					Expect(workerContainer).To(BeNil())
					Expect(found).To(BeFalse())
				})
			})

			Context("when we cannot find the worker from the container info", func() {
				BeforeEach(func() {
					fakeProvider.GetWorkerReturns(nil, false, nil)
				})

				It("returns ErrMissingWorker", func() {
					container, found, err := pool.FindContainerForIdentifier(logger, identifier)
					Expect(err).To(Equal(ErrMissingWorker))
					Expect(container).To(BeNil())
					Expect(found).To(BeFalse())
				})
			})

			Context("when looking up the worker is successful", func() {
				var fakeWorker *workerfakes.FakeWorker

				BeforeEach(func() {
					fakeWorker = new(workerfakes.FakeWorker)
					fakeProvider.GetWorkerReturns(fakeWorker, true, nil)
				})

				It("calls to lookup the container on the worker", func() {
					pool.FindContainerForIdentifier(logger, identifier)

					Expect(fakeWorker.FindContainerByHandleCallCount()).To(Equal(1))

					_, handleArg, _ := fakeWorker.FindContainerByHandleArgsForCall(0)
					Expect(handleArg).To(Equal("some-container-handle"))
				})

				Context("when looking up the container contains an error", func() {
					It("returns the error", func() {
						fakeWorker.FindContainerByHandleReturns(nil, false, errors.New("disaster"))

						container, found, err := pool.FindContainerForIdentifier(logger, identifier)
						Expect(err).To(HaveOccurred())
						Expect(container).To(BeNil())
						Expect(found).To(BeFalse())
					})
				})

				Context("when finding the container on the worker is successful", func() {
					var (
						fakeContainer  *workerfakes.FakeContainer
						foundContainer Container
						found          bool
						err            error
					)

					BeforeEach(func() {
						fakeContainer = new(workerfakes.FakeContainer)
						fakeWorker.FindContainerByHandleReturns(fakeContainer, true, nil)
					})

					JustBeforeEach(func() {
						foundContainer, found, err = pool.FindContainerForIdentifier(logger, identifier)
					})

					It("returns the container", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(found).To(BeTrue())
						Expect(foundContainer).To(Equal(fakeContainer))
					})
				})
			})
		})
	})
})
