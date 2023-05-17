BUILD=$(shell ./scripts/build.sh)
CLEAN=$(shell ./scripts/clean.sh)
MASTER=$(shell ./scripts/master.sh)
WORKER=$(shell ./scripts/worker.sh)

TESTPOD=$(shell ./scripts/testPod.sh)
TESTSERVICE=$(shell ./scripts/testService.sh)
TESTGATEWAY=$(shell ./scripts/testGateway.sh)
TESTRS=$(shell ./scripts/testRS.sh)
TESTSTRESS=$(shell ./scripts/testStress.sh)
TESTGPU=$(shell ./scripts/testGpu.sh)
TESTSERVERLESS=$(shell ./scripts/testServerless.sh)

# Build k8s but not run
.PHONY:build
build:
	@echo $(BUILD)

# Build k8s and run as master and worker
.PHONY:master
master:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(MASTER)
	@echo $(WORKER)

# Build k8s and run as worker
.PHONY:worker
worker:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(WORKER)

# Build k8s, run a master, a worker, a Pod on one host
.PHONY:testPod
testPod: master
	@echo $(TESTPOD)

# Build k8s, run a master, a worker, a Pod, a Service on one host
.PHONY:testService
testService: master
	@echo $(TESTSERVICE)

# Build k8s, run a master, a worker, a Pod, a Service, a Gateway on one host
.PHONY:testGateway
testGateway: master
	@echo $(TESTGATEWAY)

# Build k8s, run a master, a worker, a RS, a Service on one host
.PHONY:testRS
testRS: master
	@echo $(TESTRS)

# Build k8s, run a master, a worker, a Stress Pod on one host
.PHONY:testStress
testStress: master
	@echo $(TESTSTRESS)

# Build k8s, run a master, a worker, a Gpu Job on one host
.PHONY:testGpu
testGpu: master
	@echo $(TESTGPU)

# Build k8s, run a master, a worker, a serverless test on one host
.PHONY:testServerless
testServerless: master
	@echo $(TESTSERVERLESS)

# Stop k8s and clear k8s states
.PHONY:clean
clean:
	@echo $(CLEAN)
