#!/bin/bash

set -euo pipefail

tail -n 0 -q -F /opt/blockchain_node/log/*.log* >> /proc/1/fd/1 &

/opt/blockchain_node/bin/blockchain_node daemon
/app/rosetta-helium &
node /app/helium-constructor/public/index.js
