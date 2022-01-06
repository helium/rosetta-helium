#!/bin/bash

set -euo pipefail

cat /opt/blockchain_node/config/sys.config | grep -oP '(?<=\{blessed_snapshot_block_height\, ).*?(?=\})' > /app/lbs.txt &

/opt/blockchain_node/bin/blockchain_node foreground &
/app/rosetta-helium &
node /app/helium-constructor/public/index.js
