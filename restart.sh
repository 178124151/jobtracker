#!/bin/bash
# JobTracker Restart Script (Linux)

cd "$(dirname "$0")"

echo "Restarting JobTracker..."
docker compose restart
echo "Done!"