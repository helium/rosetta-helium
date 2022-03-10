#!/bin/bash

set -euo pipefail

/opt/blockchain_node/bin/blockchain_node foreground &
/app/rosetta-helium --testnet --data="/data" &
NETWORK=testnet node /app/helium-constructor/public/index.js
