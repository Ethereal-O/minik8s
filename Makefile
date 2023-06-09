BUILD=$(shell ./scripts/build.sh)
CLEAN=$(shell ./scripts/clean.sh)
MASTER=$(shell ./scripts/master.sh)
WORKER=$(shell ./scripts/worker.sh)
WORKER_PREPARE=$(shell ./scripts/worker_prepare.sh)

TESTPOD=$(shell ./scripts/testPod.sh)
TESTSERVICE=$(shell ./scripts/testService.sh)
TESTGATEWAY=$(shell ./scripts/testGateway.sh)
TESTRS=$(shell ./scripts/testRS.sh)
TESTSTRESS=$(shell ./scripts/testStress.sh)
TESTGPU=$(shell ./scripts/testGpu.sh)
TESTSERVERLESS=$(shell ./scripts/testServerless.sh)
TESTCONSISTENCY=$(shell ./scripts/testConsistency.sh)
TESTMICROSERVICE=$(shell ./scripts/testMicroService.sh)

# Build k8s but not run
.PHONY:build
build:
	@echo $(BUILD)

# Build k8s and run as master and worker
.PHONY:master
master:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(WORKER_PREPARE)
	@echo $(MASTER)
	@echo $(WORKER)

# Build k8s and run as worker
.PHONY:worker
worker:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(WORKER_PREPARE)
	@echo $(WORKER)

# Build k8s, run a master, a worker, a Pod on one host
.PHONY:testPod
testPod:
	@echo $(TESTPOD)

# Build k8s, run a master, a worker, a Pod, a Service on one host
.PHONY:testService
testService:
	@echo $(TESTSERVICE)

# Build k8s, run a master, a worker, a Pod, a Service, a Gateway on one host
.PHONY:testGateway
testGateway:
	@echo $(TESTGATEWAY)

.PHONY:testConsistency
testConsistency:
	@echo $(TESTCONSISTENCY)

# Build k8s, run a master, a worker, a RS, a Service on one host
.PHONY:testRS
testRS:
	@echo $(TESTRS)

# Build k8s, run a master, a worker, a Stress Pod on one host
.PHONY:testStress
testStress:
	@echo $(TESTSTRESS)

# Build k8s, run a master, a worker, a Gpu Job on one host
.PHONY:testGpu
testGpu:
	@echo $(TESTGPU)

# Build k8s, run a master, a worker, a serverless test on one host
.PHONY:testServerless
testServerless:
	@echo $(TESTSERVERLESS)

# Build k8s, run a master, a worker, a microservice test on one host
.PHONY:testMicroService
testMicroService:
	@echo $(TESTMICROSERVICE)

# Stop k8s and clear k8s states
.PHONY:clean
clean:
	@echo $(CLEAN)
