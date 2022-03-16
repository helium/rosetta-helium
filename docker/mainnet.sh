#!/bin/bash

set -euo pipefail

echo '1267922' > /app/lbs.txt &

/opt/blockchain_node/bin/blockchain_node foreground &
/app/rosetta-helium --data="/data" &
node /app/helium-constructor/public/index.js
