# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a wrapper around whisper.cpp for audio transcription, providing a Makefile-based workflow to convert MP3 files to various transcript formats (TXT, VTT, SRT, LRC). The project uses a two-tier Makefile architecture for dependency management and transcription processing.

## Architecture

### Dual Makefile System
- **`Makefile`**: Main workflow for audio processing and transcription
- **`dependency.mk`**: Manages whisper.cpp dependency, model downloads, and builds

### Processing Pipeline
```
dist/mp3/*.mp3 → dist/wav/*.wav → dist/transcript/*.{txt,vtt,srt,lrc}
```

1. **Audio Conversion**: MP3 files are converted to 16kHz mono WAV format using ffmpeg
2. **Transcription**: WAV files are processed by whisper-cli with quality-optimized settings
3. **Multiple Output Formats**: Supports plain text, WebVTT, SRT, and LRC subtitle formats

### Directory Structure
- `dist/mp3/` - Source MP3 files
- `dist/wav/` - Converted WAV files (16kHz mono)
- `dist/transcript/` - Generated transcripts in various formats
- `dependency/ggml-org/whisper.cpp/` - Cloned whisper.cpp repository
- `bin/whisper-cli` - Symlink to built whisper-cli binary

## Essential Commands

### Setup and Dependencies
```bash
# Initialize project directories and dependencies
make depend

# Setup project directories only
make setup

# Clean up everything
make independ

# Remove generated files only
make teardown
```

### Transcription Workflow
```bash
# Convert MP3 to WAV (prerequisite for transcription)
make dist/wav/filename.mp3.wav

# Generate transcripts in different formats
make dist/transcript/filename.wav.txt    # Plain text
make dist/transcript/filename.wav.vtt    # WebVTT subtitles
make dist/transcript/filename.wav.srt    # SRT subtitles  
make dist/transcript/filename.wav.lrc    # LRC lyrics format
```

### Dependency Management
```bash
# Download whisper models
make -f dependency.mk download

# Build whisper-cli from source
make -f dependency.mk build

# Clean build artifacts
make -f dependency.mk clean
```

## Configuration Variables

The Makefile includes quality-optimized transcription settings:

- `MODEL`: Whisper model to use (default: `ggml-large-v3-turbo.bin`)
- `LANGUAGE`: Language setting (default: `auto`)
- `BEST_OF`: Number of best candidates (default: `8`)  
- `BEAM_SIZE`: Beam search size (default: `8`)
- `TEMPERATURE`: Sampling temperature (default: `0.1`)
- `WORD_THRESHOLD`: Word confidence threshold (default: `0.02`)
- `MAX_CONTEXT`: Maximum context tokens (default: `224`)
- `THREADS`: Number of processing threads (default: `8`)

## Quality Features

The transcription pipeline includes several quality improvements:
- **Beam search optimization**: Uses maximum supported beam size (8)
- **Non-speech token suppression**: `--suppress-nst` flag reduces hallucinations
- **Flash attention**: `--flash-attn` for improved processing
- **Word confidence filtering**: Higher threshold filters uncertain transcriptions
- **Optimized temperature**: Low temperature (0.1) for more deterministic output

## File Naming Convention

Target files follow the pattern: `dist/transcript/[source-filename].wav.[format]`

Example: `dist/mp3/podcast-episode.mp3` produces:
- `dist/wav/podcast-episode.mp3.wav`
- `dist/transcript/podcast-episode.mp3.wav.txt`
- `dist/transcript/podcast-episode.mp3.wav.srt`

## Dependencies

- **ffmpeg**: Required for MP3 to WAV conversion
- **cmake**: Required to build whisper.cpp
- **git**: Required to clone whisper.cpp repository
- **whisper.cpp**: Automatically cloned and built via dependency.mk

The project automatically handles whisper.cpp setup, including model downloads for base.en, large-v3-turbo, and large-v3-turbo-q5_0 variants.