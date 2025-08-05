import json
import os
import random
import string
from pathlib import Path
from datetime import datetime

# --- CONFIGURATION ---
OUTPUT_DIR = "./test"
DEPTH = 4  # max levels of nested objects
PROPERTIES = 3  # approx number of properties per log
COUNT = 10  # number of JSON files to generate
MIXED_ARRAYS = False  # whether arrays can have mixed value types
# ----------------------

def random_string(length=6):
    return ''.join(random.choices(string.ascii_lowercase, k=length))

def random_primitive():
    return random.choice([
        random.randint(0, 1000),
        random.random(),
        random_string(),
        True,
        False,
        None
    ])

def random_array(depth, mixed):
    length = random.randint(2, 5)
    if mixed:
        return [random_value(depth - 1, mixed) for _ in range(length)]
    else:
        value = random_value(depth - 1, mixed)
        return [value for _ in range(length)]

def random_object(depth, remaining_props, mixed):
    obj = {}
    num_props = random.randint(1, remaining_props)
    for _ in range(num_props):
        key = random_string()
        obj[key] = random_value(depth - 1, mixed)
    return obj

def random_value(depth, mixed):
    if depth <= 0:
        return random_primitive()
    
    choice = random.choice(['primitive', 'array', 'object'])
    if choice == 'primitive':
        return random_primitive()
    elif choice == 'array':
        return random_array(depth, mixed)
    else:
        return random_object(depth, max(1, depth * 2), mixed)

def generate_json_log(depth, properties, mixed):
    log = {
        "timestamp": datetime.utcnow().isoformat() + "Z",
        "level": random.choice(["DEBUG", "INFO", "WARN", "ERROR"]),
        "service": f"{random_string(5)}-service",
        "trace_id": random_string(32),
        "payload": {}
    }
    
    remaining_props = properties
    while remaining_props > 0:
        key = random_string()
        value = random_value(depth, mixed)
        log["payload"][key] = value
        remaining_props -= 1
    
    return log

def main():
    output_path = Path(OUTPUT_DIR)
    output_path.mkdir(parents=True, exist_ok=True)

    for i in range(COUNT):
        log_data = generate_json_log(DEPTH, PROPERTIES, MIXED_ARRAYS)
        filename = f"test{'' if i == 0 else i}.json"
        file_path = output_path / filename
        with open(file_path, "w") as f:
            json.dump(log_data, f, indent=2)
        print(f"Generated {file_path}")

if __name__ == "__main__":
    main()
