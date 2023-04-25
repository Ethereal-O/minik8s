BUILD=$(shell ./scripts/build.sh)
TESTPOD=$(shell ./scripts/testPod.sh)
TESTSERVICE=$(shell ./scripts/testService.sh)
CLEAN=$(shell ./scripts/clean.sh)
STOP=$(shell ./scripts/stop.sh)
MASTER=$(shell ./scripts/master.sh)
WORKER=$(shell ./scripts/worker.sh)

# Build k8s but not run
build:
	@echo $(BUILD)
# Build k8s, run a master, a worker and a Pod on one host
testPod:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(MASTER)
	@echo $(WORKER)
	@echo $(TESTPOD)
# Build k8s, run a master, a worker, a Pod and a Service on one host
testService:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(MASTER)
	@echo $(WORKER)
	@echo $(TESTSERVICE)
# Build k8s and run as master and worker
master:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(MASTER)
	@echo $(WORKER)
# Build k8s and run as worker
worker:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(WORKER)
# Stop k8s and clear etcd and containers
clean:
	@echo $(CLEAN)
# Stop everything, include k8s and environment
stop:
	@echo $(STOP)

.PHONY:build
