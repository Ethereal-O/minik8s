# note: call scripts from /scripts
BUILD=$(shell ./scripts/build.sh)
TESTRR=$(shell ./scripts/testRR.sh)
CLEAN=$(shell ./scripts/clean.sh)
RUN=$(shell ./scripts/run.sh)
build:
	@echo $(CLEAN)
	@echo $(BUILD)
run:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(RUN)
testRR:
	@echo $(CLEAN)
	@echo $(BUILD)
	@echo $(RUN)
	@echo $(TESTRR)
clean:
	@echo $(CLEAN)
.PHONY:build
