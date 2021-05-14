**THIS IS NOT PRODUCTION READY. USE AT YOUR OWN RISK.**

## Overview
Bare bones Rosetta API implementation of the Helium `blockchain-node`

## Quick setup

#### Build container
```text
docker build . -t rosetta-helium:latest
```

#### Run container
Local data is stored in `helium-data`
```text
docker run -d --rm --ulimit "nofile=100000:100000" -v "$(pwd)/helium-data:/app/blockchain-node/_build/dev/rel/blockchain_node/data" -p 8080:8080 -p 44158:44158 -p 4467:4467 rosetta-helium:latest
```

#### Rosetta CLI check
```text
rosetta-cli check:data --configuration-file rosetta-cli-config/mainnet/config.json
```
(Please wait a few minutes for the Helium node to initialize before running this command)

## Implemented currencies
- HNT
- HST

### Unimplemented currencies
- DC (DCs are not implemented as they cannot be actively traded)

## Data API transactions
Transactions support for reading from the Data API

### Implemented
`payment_v1`

`payment_v2`

`reward_v1`

TODO: `reward_v2`

`security_coinbase_v1`

TODO: `security_exchange_v1`

TODO: `token_burn_v1`

TODO: `transfer_hotspot_v1`

TODO: `create_htlc_v1`*

TODO: `redeem_htlc_v1`*

### Fee-only transactions (Only recording implicit_burns for HNT deductions)
`add_gateway_v1`

TODO: `assert_location_v1`

TODO: `assert_location_v2`

TODO: `oui_v1`

TODO: `routing_v1`

TODO: `state_channel_open_v1`

### Unimplemented transactions
`dc_coinbase_v1`

`state_channel_close_v1`

## Construction API transactions
Transaction support for creation via the construction API

### Implemented
TODO: `payment_v1`

TODO: `payment_v2`

TODO: `security_exchange_v1`

TODO: `create_htlc_v1`*

TODO: `redeem_htlc_v1`*


*Lower priority 

#### Buy me an HNT coffee :)
HNT Address: 13Ey7fZfdQB7C8FRiWfuNKshA8si7wgHVeMdGbdCkM5gyry7G88
