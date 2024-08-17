#!/bin/bash

./bin/veritas &
sleep 2
PID=$!

./bin/schema-validator
EXIT_CODE=$?

echo "Sending shutdown signal to process $PID..."
kill -SIGINT "$PID"
sleep 2

if curl -s http://localhost:8080/ > /dev/null; then
  echo "Error: The server is still running."
  exit 1
else
  echo "Server has been shut down."
fi

exit "$EXIT_CODE"
