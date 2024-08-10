import sys
import resource
import time
import json

def main():
    # The script to execute and the input to pass
    script, input_data = sys.argv[1], sys.argv[2]

    # Start measuring time and memory usage
    start_time = time.time()
    start_mem = resource.getrusage(resource.RUSAGE_CHILDREN).ru_maxrss

    # Execute the script with the input
    exec(open(script).read())

    # Measure time and memory usage after execution
    end_time = time.time()
    end_mem = resource.getrusage(resource.RUSAGE_CHILDREN).ru_maxrss

    # Calculate runtime and memory used
    runtime = end_time - start_time
    memory_used = (end_mem - start_mem) / 1024  # Convert to MB

    # Output the result
    result = {
        "output": output,  # Replace with actual output from the code
        "runtime": runtime,
        "memory_used": memory_used
    }
    
    print(json.dumps(result))

if __name__ == "__main__":
    main()
