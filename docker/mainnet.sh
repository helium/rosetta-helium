#!/bin/bash

set -euo pipefail

cat /opt/blockchain_node/config/sys.config | grep -oP '(?<=\{blessed_snapshot_block_height\, ).*?(?=\})' > /app/lbs.txt &
tail -n 0 -q -F /opt/blockchain_node/log/*.log* >> /proc/1/fd/1 &

/opt/blockchain_node/bin/blockchain_node daemon
/app/rosetta-helium &
node /app/helium-constructor/public/index.js
