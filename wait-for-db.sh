#!/bin/bash
host="$1"
port="$2"

echo "Check database connection at $host:$port"
counter=0

while ! (echo > /dev/tcp/$host/$port) >/dev/null 2>&1; do
  [ $counter -eq 0 ] && echo "Waiting for $host:$port to become available..."
  counter=$((counter+1))
  sleep 1
done

# Sleep few more seconds to make sure that db is fully functional after the tcp port is open.
sleep 4

[ $counter -gt 0 ] && echo "Database available after $counter seconds"
echo "Database $host:$port ready!"
