#! /bin/bash

/app/blockchain-node/bin/blockchain_node foreground& /app/rosetta-helium --testnet& NETWORK=testnet node /app/helium-constructor/public/index.js;