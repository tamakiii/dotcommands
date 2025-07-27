import subprocess
import json
import time
import os

def query_ollama(model, prompt):
    start_time = time.time()
    result = subprocess.run([
        'ollama', 'run', model, prompt
    ], capture_output=True, text=True)
    end_time = time.time()
    
    return {
        'response': result.stdout.strip(),
        'time': end_time - start_time,
        'model': model
    }

# Test prompts for your use cases
test_cases = [
    # Scoring task
    "Rate this code quality from 1-10 and explain: def calc(x): return x*2+1",
    
    # Reasoning task  
    "A function runs every hour and checks a log file. If it finds 'ERROR', it should send an alert. Write the logic steps.",
    
    # Structured output
    "Extract key information from this agent conversation and format as JSON: 'User asked to implement login. Agent created auth.py with bcrypt hashing.'"
]

def load_preferred_models():
    """Load models from preferred-models.json"""
    script_dir = os.path.dirname(os.path.abspath(__file__))
    json_path = os.path.join(script_dir, 'preferred-models.json')
    
    with open(json_path, 'r') as f:
        data = json.load(f)
    
    return [model['name'] for model in data['models']]

models = load_preferred_models()

for prompt in test_cases:
    print(f"\n=== Testing: {prompt[:50]}... ===")
    for model in models:
        result = query_ollama(model, prompt)
        print(f"\n{model} ({result['time']:.1f}s):")
        print(result['response'][:200] + "...")
