# note: call scripts from /scripts
BUILD=$(shell ./scripts/build.sh)
build:
	@echo $(BUILD)
.PHONY:build
