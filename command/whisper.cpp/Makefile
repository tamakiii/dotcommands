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
LANGUAGE := auto

# Quality-focused transcription settings
BEST_OF := 8
BEAM_SIZE := 8
TEMPERATURE := 0.1
WORD_THRESHOLD := 0.02
MAX_CONTEXT := 224
THREADS := 8

dist/transcript/%.wav.txt: | dependency/ggml-org/whisper.cpp
	bin/whisper-cli \
		--model $|/models/$(MODEL) \
		--file dist/wav/$(patsubst dist/transcript/%.txt,%,$@) \
		--language $(LANGUAGE) \
		--best-of $(BEST_OF) \
		--beam-size $(BEAM_SIZE) \
		--temperature $(TEMPERATURE) \
		--word-thold $(WORD_THRESHOLD) \
		--max-context $(MAX_CONTEXT) \
		--threads $(THREADS) \
		--suppress-nst \
		--flash-attn \
		--print-progress \
		> $@
