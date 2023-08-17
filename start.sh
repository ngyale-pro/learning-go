#!/bin/sh

# Automatically exit if a command fails
set -e 


echo "---------- Run DB migration ----------"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up
echo "---------- Start application ----------"
exec "$@" # Takes all parameter as input and execute it

