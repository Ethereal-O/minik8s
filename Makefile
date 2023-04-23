# note: call scripts from /scripts
BUILD=$(shell ./scripts/build.sh)
TESTPOD=$(shell ./scripts/testPod.sh)
CLEAN=$(shell ./scripts/clean.sh)
RUN=$(shell ./scripts/run.sh)
build:
	@echo $(CLEAN)
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
.PHONY:build
