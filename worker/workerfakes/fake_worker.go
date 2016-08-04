// This file was generated by counterfeiter
package workerfakes

import (
	"os"
	"sync"
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/worker"
	"code.cloudfoundry.org/lager"
)

type FakeWorker struct {
	CreateContainerStub        func(lager.Logger, <-chan os.Signal, worker.ImageFetchingDelegate, worker.Identifier, worker.Metadata, worker.ContainerSpec, atc.ResourceTypes) (worker.Container, error)
	createContainerMutex       sync.RWMutex
	createContainerArgsForCall []struct {
		arg1 lager.Logger
		arg2 <-chan os.Signal
		arg3 worker.ImageFetchingDelegate
		arg4 worker.Identifier
		arg5 worker.Metadata
		arg6 worker.ContainerSpec
		arg7 atc.ResourceTypes
	}
	createContainerReturns struct {
		result1 worker.Container
		result2 error
	}
	FindContainerForIdentifierStub        func(lager.Logger, worker.Identifier) (worker.Container, bool, error)
	findContainerForIdentifierMutex       sync.RWMutex
	findContainerForIdentifierArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.Identifier
	}
	findContainerForIdentifierReturns struct {
		result1 worker.Container
		result2 bool
		result3 error
	}
	LookupContainerStub        func(lager.Logger, string) (worker.Container, bool, error)
	lookupContainerMutex       sync.RWMutex
	lookupContainerArgsForCall []struct {
		arg1 lager.Logger
		arg2 string
	}
	lookupContainerReturns struct {
		result1 worker.Container
		result2 bool
		result3 error
	}
	FindResourceTypeByPathStub        func(path string) (atc.WorkerResourceType, bool)
	findResourceTypeByPathMutex       sync.RWMutex
	findResourceTypeByPathArgsForCall []struct {
		path string
	}
	findResourceTypeByPathReturns struct {
		result1 atc.WorkerResourceType
		result2 bool
	}
	FindVolumeStub        func(lager.Logger, worker.VolumeSpec) (worker.Volume, bool, error)
	findVolumeMutex       sync.RWMutex
	findVolumeArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.VolumeSpec
	}
	findVolumeReturns struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}
	CreateVolumeStub        func(logger lager.Logger, vs worker.VolumeSpec, teamID int) (worker.Volume, error)
	createVolumeMutex       sync.RWMutex
	createVolumeArgsForCall []struct {
		logger lager.Logger
		vs     worker.VolumeSpec
		teamID int
	}
	createVolumeReturns struct {
		result1 worker.Volume
		result2 error
	}
	ListVolumesStub        func(lager.Logger, worker.VolumeProperties) ([]worker.Volume, error)
	listVolumesMutex       sync.RWMutex
	listVolumesArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.VolumeProperties
	}
	listVolumesReturns struct {
		result1 []worker.Volume
		result2 error
	}
	LookupVolumeStub        func(lager.Logger, string) (worker.Volume, bool, error)
	lookupVolumeMutex       sync.RWMutex
	lookupVolumeArgsForCall []struct {
		arg1 lager.Logger
		arg2 string
	}
	lookupVolumeReturns struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}
	SatisfyingStub        func(worker.WorkerSpec, atc.ResourceTypes) (worker.Worker, error)
	satisfyingMutex       sync.RWMutex
	satisfyingArgsForCall []struct {
		arg1 worker.WorkerSpec
		arg2 atc.ResourceTypes
	}
	satisfyingReturns struct {
		result1 worker.Worker
		result2 error
	}
	AllSatisfyingStub        func(worker.WorkerSpec, atc.ResourceTypes) ([]worker.Worker, error)
	allSatisfyingMutex       sync.RWMutex
	allSatisfyingArgsForCall []struct {
		arg1 worker.WorkerSpec
		arg2 atc.ResourceTypes
	}
	allSatisfyingReturns struct {
		result1 []worker.Worker
		result2 error
	}
	GetWorkerStub        func(workerName string) (worker.Worker, error)
	getWorkerMutex       sync.RWMutex
	getWorkerArgsForCall []struct {
		workerName string
	}
	getWorkerReturns struct {
		result1 worker.Worker
		result2 error
	}
	ActiveContainersStub        func() int
	activeContainersMutex       sync.RWMutex
	activeContainersArgsForCall []struct{}
	activeContainersReturns     struct {
		result1 int
	}
	DescriptionStub        func() string
	descriptionMutex       sync.RWMutex
	descriptionArgsForCall []struct{}
	descriptionReturns     struct {
		result1 string
	}
	NameStub        func() string
	nameMutex       sync.RWMutex
	nameArgsForCall []struct{}
	nameReturns     struct {
		result1 string
	}
	UptimeStub        func() time.Duration
	uptimeMutex       sync.RWMutex
	uptimeArgsForCall []struct{}
	uptimeReturns     struct {
		result1 time.Duration
	}
	IsOwnedByTeamStub        func() bool
	isOwnedByTeamMutex       sync.RWMutex
	isOwnedByTeamArgsForCall []struct{}
	isOwnedByTeamReturns     struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeWorker) CreateContainer(arg1 lager.Logger, arg2 <-chan os.Signal, arg3 worker.ImageFetchingDelegate, arg4 worker.Identifier, arg5 worker.Metadata, arg6 worker.ContainerSpec, arg7 atc.ResourceTypes) (worker.Container, error) {
	fake.createContainerMutex.Lock()
	fake.createContainerArgsForCall = append(fake.createContainerArgsForCall, struct {
		arg1 lager.Logger
		arg2 <-chan os.Signal
		arg3 worker.ImageFetchingDelegate
		arg4 worker.Identifier
		arg5 worker.Metadata
		arg6 worker.ContainerSpec
		arg7 atc.ResourceTypes
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7})
	fake.recordInvocation("CreateContainer", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7})
	fake.createContainerMutex.Unlock()
	if fake.CreateContainerStub != nil {
		return fake.CreateContainerStub(arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	} else {
		return fake.createContainerReturns.result1, fake.createContainerReturns.result2
	}
}

func (fake *FakeWorker) CreateContainerCallCount() int {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	return len(fake.createContainerArgsForCall)
}

func (fake *FakeWorker) CreateContainerArgsForCall(i int) (lager.Logger, <-chan os.Signal, worker.ImageFetchingDelegate, worker.Identifier, worker.Metadata, worker.ContainerSpec, atc.ResourceTypes) {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	return fake.createContainerArgsForCall[i].arg1, fake.createContainerArgsForCall[i].arg2, fake.createContainerArgsForCall[i].arg3, fake.createContainerArgsForCall[i].arg4, fake.createContainerArgsForCall[i].arg5, fake.createContainerArgsForCall[i].arg6, fake.createContainerArgsForCall[i].arg7
}

func (fake *FakeWorker) CreateContainerReturns(result1 worker.Container, result2 error) {
	fake.CreateContainerStub = nil
	fake.createContainerReturns = struct {
		result1 worker.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeWorker) FindContainerForIdentifier(arg1 lager.Logger, arg2 worker.Identifier) (worker.Container, bool, error) {
	fake.findContainerForIdentifierMutex.Lock()
	fake.findContainerForIdentifierArgsForCall = append(fake.findContainerForIdentifierArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.Identifier
	}{arg1, arg2})
	fake.recordInvocation("FindContainerForIdentifier", []interface{}{arg1, arg2})
	fake.findContainerForIdentifierMutex.Unlock()
	if fake.FindContainerForIdentifierStub != nil {
		return fake.FindContainerForIdentifierStub(arg1, arg2)
	} else {
		return fake.findContainerForIdentifierReturns.result1, fake.findContainerForIdentifierReturns.result2, fake.findContainerForIdentifierReturns.result3
	}
}

func (fake *FakeWorker) FindContainerForIdentifierCallCount() int {
	fake.findContainerForIdentifierMutex.RLock()
	defer fake.findContainerForIdentifierMutex.RUnlock()
	return len(fake.findContainerForIdentifierArgsForCall)
}

func (fake *FakeWorker) FindContainerForIdentifierArgsForCall(i int) (lager.Logger, worker.Identifier) {
	fake.findContainerForIdentifierMutex.RLock()
	defer fake.findContainerForIdentifierMutex.RUnlock()
	return fake.findContainerForIdentifierArgsForCall[i].arg1, fake.findContainerForIdentifierArgsForCall[i].arg2
}

func (fake *FakeWorker) FindContainerForIdentifierReturns(result1 worker.Container, result2 bool, result3 error) {
	fake.FindContainerForIdentifierStub = nil
	fake.findContainerForIdentifierReturns = struct {
		result1 worker.Container
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorker) LookupContainer(arg1 lager.Logger, arg2 string) (worker.Container, bool, error) {
	fake.lookupContainerMutex.Lock()
	fake.lookupContainerArgsForCall = append(fake.lookupContainerArgsForCall, struct {
		arg1 lager.Logger
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("LookupContainer", []interface{}{arg1, arg2})
	fake.lookupContainerMutex.Unlock()
	if fake.LookupContainerStub != nil {
		return fake.LookupContainerStub(arg1, arg2)
	} else {
		return fake.lookupContainerReturns.result1, fake.lookupContainerReturns.result2, fake.lookupContainerReturns.result3
	}
}

func (fake *FakeWorker) LookupContainerCallCount() int {
	fake.lookupContainerMutex.RLock()
	defer fake.lookupContainerMutex.RUnlock()
	return len(fake.lookupContainerArgsForCall)
}

func (fake *FakeWorker) LookupContainerArgsForCall(i int) (lager.Logger, string) {
	fake.lookupContainerMutex.RLock()
	defer fake.lookupContainerMutex.RUnlock()
	return fake.lookupContainerArgsForCall[i].arg1, fake.lookupContainerArgsForCall[i].arg2
}

func (fake *FakeWorker) LookupContainerReturns(result1 worker.Container, result2 bool, result3 error) {
	fake.LookupContainerStub = nil
	fake.lookupContainerReturns = struct {
		result1 worker.Container
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorker) FindResourceTypeByPath(path string) (atc.WorkerResourceType, bool) {
	fake.findResourceTypeByPathMutex.Lock()
	fake.findResourceTypeByPathArgsForCall = append(fake.findResourceTypeByPathArgsForCall, struct {
		path string
	}{path})
	fake.recordInvocation("FindResourceTypeByPath", []interface{}{path})
	fake.findResourceTypeByPathMutex.Unlock()
	if fake.FindResourceTypeByPathStub != nil {
		return fake.FindResourceTypeByPathStub(path)
	} else {
		return fake.findResourceTypeByPathReturns.result1, fake.findResourceTypeByPathReturns.result2
	}
}

func (fake *FakeWorker) FindResourceTypeByPathCallCount() int {
	fake.findResourceTypeByPathMutex.RLock()
	defer fake.findResourceTypeByPathMutex.RUnlock()
	return len(fake.findResourceTypeByPathArgsForCall)
}

func (fake *FakeWorker) FindResourceTypeByPathArgsForCall(i int) string {
	fake.findResourceTypeByPathMutex.RLock()
	defer fake.findResourceTypeByPathMutex.RUnlock()
	return fake.findResourceTypeByPathArgsForCall[i].path
}

func (fake *FakeWorker) FindResourceTypeByPathReturns(result1 atc.WorkerResourceType, result2 bool) {
	fake.FindResourceTypeByPathStub = nil
	fake.findResourceTypeByPathReturns = struct {
		result1 atc.WorkerResourceType
		result2 bool
	}{result1, result2}
}

func (fake *FakeWorker) FindVolume(arg1 lager.Logger, arg2 worker.VolumeSpec) (worker.Volume, bool, error) {
	fake.findVolumeMutex.Lock()
	fake.findVolumeArgsForCall = append(fake.findVolumeArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.VolumeSpec
	}{arg1, arg2})
	fake.recordInvocation("FindVolume", []interface{}{arg1, arg2})
	fake.findVolumeMutex.Unlock()
	if fake.FindVolumeStub != nil {
		return fake.FindVolumeStub(arg1, arg2)
	} else {
		return fake.findVolumeReturns.result1, fake.findVolumeReturns.result2, fake.findVolumeReturns.result3
	}
}

func (fake *FakeWorker) FindVolumeCallCount() int {
	fake.findVolumeMutex.RLock()
	defer fake.findVolumeMutex.RUnlock()
	return len(fake.findVolumeArgsForCall)
}

func (fake *FakeWorker) FindVolumeArgsForCall(i int) (lager.Logger, worker.VolumeSpec) {
	fake.findVolumeMutex.RLock()
	defer fake.findVolumeMutex.RUnlock()
	return fake.findVolumeArgsForCall[i].arg1, fake.findVolumeArgsForCall[i].arg2
}

func (fake *FakeWorker) FindVolumeReturns(result1 worker.Volume, result2 bool, result3 error) {
	fake.FindVolumeStub = nil
	fake.findVolumeReturns = struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorker) CreateVolume(logger lager.Logger, vs worker.VolumeSpec, teamID int) (worker.Volume, error) {
	fake.createVolumeMutex.Lock()
	fake.createVolumeArgsForCall = append(fake.createVolumeArgsForCall, struct {
		logger lager.Logger
		vs     worker.VolumeSpec
		teamID int
	}{logger, vs, teamID})
	fake.recordInvocation("CreateVolume", []interface{}{logger, vs, teamID})
	fake.createVolumeMutex.Unlock()
	if fake.CreateVolumeStub != nil {
		return fake.CreateVolumeStub(logger, vs, teamID)
	} else {
		return fake.createVolumeReturns.result1, fake.createVolumeReturns.result2
	}
}

func (fake *FakeWorker) CreateVolumeCallCount() int {
	fake.createVolumeMutex.RLock()
	defer fake.createVolumeMutex.RUnlock()
	return len(fake.createVolumeArgsForCall)
}

func (fake *FakeWorker) CreateVolumeArgsForCall(i int) (lager.Logger, worker.VolumeSpec, int) {
	fake.createVolumeMutex.RLock()
	defer fake.createVolumeMutex.RUnlock()
	return fake.createVolumeArgsForCall[i].logger, fake.createVolumeArgsForCall[i].vs, fake.createVolumeArgsForCall[i].teamID
}

func (fake *FakeWorker) CreateVolumeReturns(result1 worker.Volume, result2 error) {
	fake.CreateVolumeStub = nil
	fake.createVolumeReturns = struct {
		result1 worker.Volume
		result2 error
	}{result1, result2}
}

func (fake *FakeWorker) ListVolumes(arg1 lager.Logger, arg2 worker.VolumeProperties) ([]worker.Volume, error) {
	fake.listVolumesMutex.Lock()
	fake.listVolumesArgsForCall = append(fake.listVolumesArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.VolumeProperties
	}{arg1, arg2})
	fake.recordInvocation("ListVolumes", []interface{}{arg1, arg2})
	fake.listVolumesMutex.Unlock()
	if fake.ListVolumesStub != nil {
		return fake.ListVolumesStub(arg1, arg2)
	} else {
		return fake.listVolumesReturns.result1, fake.listVolumesReturns.result2
	}
}

func (fake *FakeWorker) ListVolumesCallCount() int {
	fake.listVolumesMutex.RLock()
	defer fake.listVolumesMutex.RUnlock()
	return len(fake.listVolumesArgsForCall)
}

func (fake *FakeWorker) ListVolumesArgsForCall(i int) (lager.Logger, worker.VolumeProperties) {
	fake.listVolumesMutex.RLock()
	defer fake.listVolumesMutex.RUnlock()
	return fake.listVolumesArgsForCall[i].arg1, fake.listVolumesArgsForCall[i].arg2
}

func (fake *FakeWorker) ListVolumesReturns(result1 []worker.Volume, result2 error) {
	fake.ListVolumesStub = nil
	fake.listVolumesReturns = struct {
		result1 []worker.Volume
		result2 error
	}{result1, result2}
}

func (fake *FakeWorker) LookupVolume(arg1 lager.Logger, arg2 string) (worker.Volume, bool, error) {
	fake.lookupVolumeMutex.Lock()
	fake.lookupVolumeArgsForCall = append(fake.lookupVolumeArgsForCall, struct {
		arg1 lager.Logger
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("LookupVolume", []interface{}{arg1, arg2})
	fake.lookupVolumeMutex.Unlock()
	if fake.LookupVolumeStub != nil {
		return fake.LookupVolumeStub(arg1, arg2)
	} else {
		return fake.lookupVolumeReturns.result1, fake.lookupVolumeReturns.result2, fake.lookupVolumeReturns.result3
	}
}

func (fake *FakeWorker) LookupVolumeCallCount() int {
	fake.lookupVolumeMutex.RLock()
	defer fake.lookupVolumeMutex.RUnlock()
	return len(fake.lookupVolumeArgsForCall)
}

func (fake *FakeWorker) LookupVolumeArgsForCall(i int) (lager.Logger, string) {
	fake.lookupVolumeMutex.RLock()
	defer fake.lookupVolumeMutex.RUnlock()
	return fake.lookupVolumeArgsForCall[i].arg1, fake.lookupVolumeArgsForCall[i].arg2
}

func (fake *FakeWorker) LookupVolumeReturns(result1 worker.Volume, result2 bool, result3 error) {
	fake.LookupVolumeStub = nil
	fake.lookupVolumeReturns = struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeWorker) Satisfying(arg1 worker.WorkerSpec, arg2 atc.ResourceTypes) (worker.Worker, error) {
	fake.satisfyingMutex.Lock()
	fake.satisfyingArgsForCall = append(fake.satisfyingArgsForCall, struct {
		arg1 worker.WorkerSpec
		arg2 atc.ResourceTypes
	}{arg1, arg2})
	fake.recordInvocation("Satisfying", []interface{}{arg1, arg2})
	fake.satisfyingMutex.Unlock()
	if fake.SatisfyingStub != nil {
		return fake.SatisfyingStub(arg1, arg2)
	} else {
		return fake.satisfyingReturns.result1, fake.satisfyingReturns.result2
	}
}

func (fake *FakeWorker) SatisfyingCallCount() int {
	fake.satisfyingMutex.RLock()
	defer fake.satisfyingMutex.RUnlock()
	return len(fake.satisfyingArgsForCall)
}

func (fake *FakeWorker) SatisfyingArgsForCall(i int) (worker.WorkerSpec, atc.ResourceTypes) {
	fake.satisfyingMutex.RLock()
	defer fake.satisfyingMutex.RUnlock()
	return fake.satisfyingArgsForCall[i].arg1, fake.satisfyingArgsForCall[i].arg2
}

func (fake *FakeWorker) SatisfyingReturns(result1 worker.Worker, result2 error) {
	fake.SatisfyingStub = nil
	fake.satisfyingReturns = struct {
		result1 worker.Worker
		result2 error
	}{result1, result2}
}

func (fake *FakeWorker) AllSatisfying(arg1 worker.WorkerSpec, arg2 atc.ResourceTypes) ([]worker.Worker, error) {
	fake.allSatisfyingMutex.Lock()
	fake.allSatisfyingArgsForCall = append(fake.allSatisfyingArgsForCall, struct {
		arg1 worker.WorkerSpec
		arg2 atc.ResourceTypes
	}{arg1, arg2})
	fake.recordInvocation("AllSatisfying", []interface{}{arg1, arg2})
	fake.allSatisfyingMutex.Unlock()
	if fake.AllSatisfyingStub != nil {
		return fake.AllSatisfyingStub(arg1, arg2)
	} else {
		return fake.allSatisfyingReturns.result1, fake.allSatisfyingReturns.result2
	}
}

func (fake *FakeWorker) AllSatisfyingCallCount() int {
	fake.allSatisfyingMutex.RLock()
	defer fake.allSatisfyingMutex.RUnlock()
	return len(fake.allSatisfyingArgsForCall)
}

func (fake *FakeWorker) AllSatisfyingArgsForCall(i int) (worker.WorkerSpec, atc.ResourceTypes) {
	fake.allSatisfyingMutex.RLock()
	defer fake.allSatisfyingMutex.RUnlock()
	return fake.allSatisfyingArgsForCall[i].arg1, fake.allSatisfyingArgsForCall[i].arg2
}

func (fake *FakeWorker) AllSatisfyingReturns(result1 []worker.Worker, result2 error) {
	fake.AllSatisfyingStub = nil
	fake.allSatisfyingReturns = struct {
		result1 []worker.Worker
		result2 error
	}{result1, result2}
}

func (fake *FakeWorker) GetWorker(workerName string) (worker.Worker, error) {
	fake.getWorkerMutex.Lock()
	fake.getWorkerArgsForCall = append(fake.getWorkerArgsForCall, struct {
		workerName string
	}{workerName})
	fake.recordInvocation("GetWorker", []interface{}{workerName})
	fake.getWorkerMutex.Unlock()
	if fake.GetWorkerStub != nil {
		return fake.GetWorkerStub(workerName)
	} else {
		return fake.getWorkerReturns.result1, fake.getWorkerReturns.result2
	}
}

func (fake *FakeWorker) GetWorkerCallCount() int {
	fake.getWorkerMutex.RLock()
	defer fake.getWorkerMutex.RUnlock()
	return len(fake.getWorkerArgsForCall)
}

func (fake *FakeWorker) GetWorkerArgsForCall(i int) string {
	fake.getWorkerMutex.RLock()
	defer fake.getWorkerMutex.RUnlock()
	return fake.getWorkerArgsForCall[i].workerName
}

func (fake *FakeWorker) GetWorkerReturns(result1 worker.Worker, result2 error) {
	fake.GetWorkerStub = nil
	fake.getWorkerReturns = struct {
		result1 worker.Worker
		result2 error
	}{result1, result2}
}

func (fake *FakeWorker) ActiveContainers() int {
	fake.activeContainersMutex.Lock()
	fake.activeContainersArgsForCall = append(fake.activeContainersArgsForCall, struct{}{})
	fake.recordInvocation("ActiveContainers", []interface{}{})
	fake.activeContainersMutex.Unlock()
	if fake.ActiveContainersStub != nil {
		return fake.ActiveContainersStub()
	} else {
		return fake.activeContainersReturns.result1
	}
}

func (fake *FakeWorker) ActiveContainersCallCount() int {
	fake.activeContainersMutex.RLock()
	defer fake.activeContainersMutex.RUnlock()
	return len(fake.activeContainersArgsForCall)
}

func (fake *FakeWorker) ActiveContainersReturns(result1 int) {
	fake.ActiveContainersStub = nil
	fake.activeContainersReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeWorker) Description() string {
	fake.descriptionMutex.Lock()
	fake.descriptionArgsForCall = append(fake.descriptionArgsForCall, struct{}{})
	fake.recordInvocation("Description", []interface{}{})
	fake.descriptionMutex.Unlock()
	if fake.DescriptionStub != nil {
		return fake.DescriptionStub()
	} else {
		return fake.descriptionReturns.result1
	}
}

func (fake *FakeWorker) DescriptionCallCount() int {
	fake.descriptionMutex.RLock()
	defer fake.descriptionMutex.RUnlock()
	return len(fake.descriptionArgsForCall)
}

func (fake *FakeWorker) DescriptionReturns(result1 string) {
	fake.DescriptionStub = nil
	fake.descriptionReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeWorker) Name() string {
	fake.nameMutex.Lock()
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct{}{})
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub()
	} else {
		return fake.nameReturns.result1
	}
}

func (fake *FakeWorker) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *FakeWorker) NameReturns(result1 string) {
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeWorker) Uptime() time.Duration {
	fake.uptimeMutex.Lock()
	fake.uptimeArgsForCall = append(fake.uptimeArgsForCall, struct{}{})
	fake.recordInvocation("Uptime", []interface{}{})
	fake.uptimeMutex.Unlock()
	if fake.UptimeStub != nil {
		return fake.UptimeStub()
	} else {
		return fake.uptimeReturns.result1
	}
}

func (fake *FakeWorker) UptimeCallCount() int {
	fake.uptimeMutex.RLock()
	defer fake.uptimeMutex.RUnlock()
	return len(fake.uptimeArgsForCall)
}

func (fake *FakeWorker) UptimeReturns(result1 time.Duration) {
	fake.UptimeStub = nil
	fake.uptimeReturns = struct {
		result1 time.Duration
	}{result1}
}

func (fake *FakeWorker) IsOwnedByTeam() bool {
	fake.isOwnedByTeamMutex.Lock()
	fake.isOwnedByTeamArgsForCall = append(fake.isOwnedByTeamArgsForCall, struct{}{})
	fake.recordInvocation("IsOwnedByTeam", []interface{}{})
	fake.isOwnedByTeamMutex.Unlock()
	if fake.IsOwnedByTeamStub != nil {
		return fake.IsOwnedByTeamStub()
	} else {
		return fake.isOwnedByTeamReturns.result1
	}
}

func (fake *FakeWorker) IsOwnedByTeamCallCount() int {
	fake.isOwnedByTeamMutex.RLock()
	defer fake.isOwnedByTeamMutex.RUnlock()
	return len(fake.isOwnedByTeamArgsForCall)
}

func (fake *FakeWorker) IsOwnedByTeamReturns(result1 bool) {
	fake.IsOwnedByTeamStub = nil
	fake.isOwnedByTeamReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeWorker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	fake.findContainerForIdentifierMutex.RLock()
	defer fake.findContainerForIdentifierMutex.RUnlock()
	fake.lookupContainerMutex.RLock()
	defer fake.lookupContainerMutex.RUnlock()
	fake.findResourceTypeByPathMutex.RLock()
	defer fake.findResourceTypeByPathMutex.RUnlock()
	fake.findVolumeMutex.RLock()
	defer fake.findVolumeMutex.RUnlock()
	fake.createVolumeMutex.RLock()
	defer fake.createVolumeMutex.RUnlock()
	fake.listVolumesMutex.RLock()
	defer fake.listVolumesMutex.RUnlock()
	fake.lookupVolumeMutex.RLock()
	defer fake.lookupVolumeMutex.RUnlock()
	fake.satisfyingMutex.RLock()
	defer fake.satisfyingMutex.RUnlock()
	fake.allSatisfyingMutex.RLock()
	defer fake.allSatisfyingMutex.RUnlock()
	fake.getWorkerMutex.RLock()
	defer fake.getWorkerMutex.RUnlock()
	fake.activeContainersMutex.RLock()
	defer fake.activeContainersMutex.RUnlock()
	fake.descriptionMutex.RLock()
	defer fake.descriptionMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.uptimeMutex.RLock()
	defer fake.uptimeMutex.RUnlock()
	fake.isOwnedByTeamMutex.RLock()
	defer fake.isOwnedByTeamMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeWorker) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ worker.Worker = new(FakeWorker)
