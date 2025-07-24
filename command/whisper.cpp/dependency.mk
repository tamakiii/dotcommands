.PHONY: help setup teardown download download-ggml build clean

help:
	@cat $(firstword $(MAKEFILE_LIST))

setup: \
	dependency \
	dependency/ggml-org \
	dependency/ggml-org/whisper.cpp

teardown:
	rm -rf dependency

download: \
	download-ggml

download-ggml: | dependency/ggml-org/whisper.cpp
	(cd $| && sh ./models/download-ggml-model.sh base.en)

build: \
	dependency/ggml-org/whisper.cpp/build \
	dependency/ggml-org/whisper.cpp/build/bin/whisper-cli \
	bin \
	bin/whisper-cli

clean:
	rm -rf dependency/ggml-org/whisper.cpp/build/bin/whisper-cli
	rm -rf dependency/ggml-org/whisper.cpp/build

bin:
	-mkdir $@

bin/whisper-cli: dependency/ggml-org/whisper.cpp/build/bin/whisper-cli | bin
	ln -sf $$(realpath --relative-to=bin $<) $@

dependency/ggml-org/whisper.cpp/build: dependency/ggml-org/whisper.cpp/build/bin/whisper-cli | dependency/ggml-org/whisper.cpp
	(cd $| && cmake -B build)

dependency/ggml-org/whisper.cpp/build/bin/whisper-cli: | dependency/ggml-org/whisper.cpp
	(cd $| && cmake --build build -j --config Release)

dependency:
	-mkdir $@

dependency/ggml-org:
	-mkdir $@

dependency/ggml-org/whisper.cpp:
	git clone https://github.com/ggml-org/whisper.cpp.git $@
