.PHONY: help setup teardown

help:
	@cat $(firstword $(MAKEFILE_LIST))

setup: \
	dependency \
	dependency/ggml-org \
	dependency/ggml-org/whisper.cpp

teardown:
	rm -rf dependency

dependency:
	-mkdir $@

dependency/ggml-org:
	-mkdir $@

dependency/ggml-org/whisper.cpp:
	git clone https://github.com/ggml-org/whisper.cpp.git $@
