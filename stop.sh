#!/bin/bash
# JobTracker Stop Script (Linux)

cd "$(dirname "$0")"

echo "Stopping JobTracker..."
docker compose stop
echo "Done!"