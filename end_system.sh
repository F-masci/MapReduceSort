#!/bin/bash

# List of ports to check
ports=(45978 45979 45980 45981 45989 45990)

# Loop through each port
for port in "${ports[@]}"; do
    echo "Checking port $port..."

    # Get the PID for the process using the port
    pid=$(netstat -ano | grep ":$port" | awk '{print $5}')

    # Check if PID was found
    if [ -n "$pid" ]; then
        echo "Found process ID $pid using port $port"
        kill -9 $pid
        echo "Process $pid has been terminated."
    else
        echo "No process found using port $port."
    fi
done
