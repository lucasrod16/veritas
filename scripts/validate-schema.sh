#!/bin/bash

./bin/veritas &
PID=$!
sleep 2

./bin/schema-validator
EXIT_CODE=$?

kill -SIGINT "$PID"
sleep 2

if curl -s http://localhost:8080/ > /dev/null; then
  echo "Error: The server is still running. You may have to manually kill the process."
  exit "$EXIT_CODE"
else
  echo "Server has been shut down."
fi

exit "$EXIT_CODE"
