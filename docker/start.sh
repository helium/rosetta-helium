#!/bin/bash

set -euo pipefail

/opt/blockchain-node/bin/blockchain_node daemon
/app/rosetta-helium &
node /app/helium-constructor/public/index.js
