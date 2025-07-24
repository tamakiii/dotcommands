.PHONY: help depend independ setup teardown

help:
	@cat $(firstword $(MAKEFILE_LIST))

depend:
	make -f dependency.mk setup
	make -f dependency.mk build

independ:
	make -f dependency.mk clean
	make -f dependency.mk teardown

setup: \
	dist \
	dist/mp3 \
	dist/wav \
	dist/transcript

teardown:
	rm -rf dist/wav
	rm -rf dist/mp3
	rm -rf dist

dist:
	-mkdir $@

dist/mp3:
	-mkdir $@

dist/wav:
	-mkdir $@

dist/transcript:
	-mkdir $@

dist/wav/%.mp3.wav: dist/mp3/%.mp3
	ffmpeg -i $< -ar 16000 -ac 1 -c:a pcm_s16le $@

MODEL := ggml-large-v3-turbo.bin

dist/transcript/%.wav: | dependency/ggml-org/whisper.cpp
	# 1. Convert % part of $@ from dist/transcript/fukabori-131.mp3.wav to fukabori-131.mp3.wav
	# 2. Do ls to check if dist/wav/fukabori-131.mp3.wav exists
	# 3. Do ls to check if $|/models/$(MODEL) exists
	# Later: we'll finish this part, not now
	# bin/whisper-cli -m $|/models/ggml-large-v3-turbo.bin -f $@
