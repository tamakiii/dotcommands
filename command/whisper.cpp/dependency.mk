.PHONY: help setup teardown download build clean

help:
	@cat $(firstword $(MAKEFILE_LIST))

setup: \
	dependency \
	dependency/ggml-org \
	dependency/ggml-org/whisper.cpp

teardown:
	rm -rf dependency

# see: https://github.com/ggml-org/whisper.cpp/blob/master/models/README.md#available-models
download: \
	dependency/ggml-org/whisper.cpp/models/ggml-base.en.bin \
	dependency/ggml-org/whisper.cpp/models/ggml-large-v3-turbo-q5_0.bin \
	dependency/ggml-org/whisper.cpp/models/large-v3-turbo.bin

dependency/ggml-org/whisper.cpp/models/ggml-base.en.bin: | dependency/ggml-org/whisper.cpp
	(cd $| && sh ./models/download-ggml-model.sh base.en)

dependency/ggml-org/whisper.cpp/models/ggml-large-v3-turbo-q5_0.bin: | dependency/ggml-org/whisper.cpp
	(cd $| && sh ./models/download-ggml-model.sh large-v3-turbo-q5_0)

dependency/ggml-org/whisper.cpp/models/large-v3-turbo.bin: | dependency/ggml-org/whisper.cpp
	(cd $| && sh ./models/download-ggml-model.sh large-v3-turbo)

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
