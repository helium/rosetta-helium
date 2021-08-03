#! /bin/bash

/app/blockchain-node/bin/blockchain_node foreground& /app/rosetta-helium& NETWORK=testnet node /app/helium-constructor/public/index.js;