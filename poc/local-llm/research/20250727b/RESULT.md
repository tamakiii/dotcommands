# Local LLM Testing Results - 20250727b

## Overview

Testing of 7 preferred models from `preferred-models.json` across 3 use cases: code scoring, reasoning tasks, and structured output generation.

## Model Performance Analysis

### Response Speed (fastest to slowest)
1. **llama3.2:3b** - 2.6-12.9s (avg: 7.1s)
2. **phi3:3.8b** - 2.8-11.5s (avg: 6.1s) 
3. **phi4:latest** - 4.6-35.3s (avg: 19.8s)
4. **phi3.5:3.8b** - 5.7-495.2s (avg: 170.9s, inconsistent)
5. **qwen2.5vl:7b** - 3.9-26.0s (avg: 17.8s)
6. **gemma3:12b** - 6.6-68.1s (avg: 37.3s)
7. **qwen3:8b** - 22.5-54.4s (avg: 40.4s)

### Task-Specific Performance

#### Code Quality Scoring
- **Best**: gemma3:12b, phi4:latest (structured analysis, clear reasoning)
- **Good**: llama3.2:3b, phi3:3.8b (consistent 6-8 ratings)
- **Issues**: qwen3:8b (confused output), phi3.5:3.8b (formatting error)

#### Reasoning Tasks (Log Monitoring Logic)
- **Best**: qwen2.5vl:7b, gemma3:12b, phi4:latest (comprehensive step-by-step)
- **Good**: llama3.2:3b (clear function overview)
- **Slow**: phi3.5:3.8b (495s response time)

#### JSON Structure Generation
- **Best**: qwen2.5vl:7b, gemma3:12b, phi3:3.8b (clean, valid JSON)
- **Good**: llama3.2:3b, phi4:latest (valid but verbose)
- **Issues**: qwen3:8b (thinking process exposed)

## Key Findings

### Top Performers
1. **qwen2.5vl:7b** - Best balance of speed and quality
2. **gemma3:12b** - Highest quality output, acceptable speed
3. **llama3.2:3b** - Fastest, consistent quality

### Model Issues
- **phi3.5:3.8b**: Severe performance inconsistency (11s vs 495s)
- **qwen3:8b**: Exposes internal "thinking" process
- **phi4:latest**: Good quality but slower than smaller models

## Recommendations

### For Production Use
- **Primary**: qwen2.5vl:7b (6GB, multimodal, reliable)
- **Backup**: llama3.2:3b (2GB, fastest, good enough quality)

### For Development/Testing  
- **gemma3:12b** when quality matters more than speed
- **Avoid**: phi3.5:3.8b due to performance inconsistency

## Hardware Context
- **Test Environment**: MacBook Pro M3 Max
- **All models**: Running via Ollama
- **Memory usage**: Within acceptable ranges for listed sizes