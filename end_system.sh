#!/bin/bash

# List of ports to check
ports=(45978 45979 45980 45981 45989 45990)

# Loop through each port
for port in "${ports[@]}"; do
    echo "Checking port $port..."

    # Get the PID of the process using the port
    pid=$(lsof -t -i:$port)

    # Check if a PID was found
    if [ -n "$pid" ]; then
        echo "Found process ID $pid using port $port"
        # Attempt to terminate the process
        if kill -9 "$pid" 2>/dev/null; then
            echo "Process $pid has been terminated."
        else
            echo "Failed to terminate process $pid. Check permissions or process state."
        fi
    else
        echo "No process found using port $port."
    fi
done