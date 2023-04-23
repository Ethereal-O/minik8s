# Build k8s but not run
BUILD=$(shell ./scripts/build.sh)

# Build k8s and run testPod.sh
TESTPOD=$(shell ./scripts/testPod.sh)

# Stop k8s and clear etcd and containers
CLEAN=$(shell ./scripts/clean.sh)

# Build k8s and run, but no kubectl commands
RUN=$(shell ./scripts/run.sh)

# Prepare environment for running, only call it once after your VM startup!
PREPARE=$(shell ./scripts/prepare.sh)

# Stop everything, include k8s and environment
STOP=$(shell ./scripts/stop.sh)

# Start nsq servic
build:
	@echo $(BUILD)
run:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(RUN)
testPod:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(RUN)
	@echo $(TESTPOD)
clean:
	@echo $(CLEAN)
prepare:
	@echo $(PREPARE)
stop:
	@echo $(STOP)

.PHONY:build
