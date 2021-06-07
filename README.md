**THIS IS NOT PRODUCTION READY. USE AT YOUR OWN RISK.**

# Overview
Dockerized Rosetta API implementation mostly based off of [blockchain-node](https://github.com/helium/blockchain-node):
- Rosetta specs: [https://www.rosetta-api.org/](https://www.rosetta-api.org/)
- This is NOT a full node, but rather works off the latest snapshot as specified in `blockchain-node`. As a result, there is currently no support for historical balances or reconciliation.
- `blockchain-node` provides the basic blockchain that the Data API reads from
- `./helium-constructor` implements a simple Express server written in TypeScript exposing [helium-js](https://github.com/helium/helium-js) for Construction API actions (transaction construction, signing mechanisms, etc)

# Quick setup

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

# Contributing
It's annoying to spin up a docker container for every change that you want to make. So for local development, it is recommended that you run each part of the implementation separately.

### rosetta-helium
1. Install [golang](https://golang.org/doc/install) if you haven't yet.
2. At the root directory, run `go run .` to start the rosetta server at port `:8080`

### blockchain-node
1. Checkout my [custom version of blockchain-node](https://github.com/syuan100/blockchain-node/tree/syuan100-fee-differentiator) that accounts for [implicit burn](https://docs.helium.com/blockchain/transaction-fees/) events.
2. Run `make && make release PROFILE=devib` to build a release
3. Run `make start PROFILE=devib` to start blockchain-node at port `:4467`

### helium-constructor
1. Install `node`. I prefer [nvm](https://github.com/nvm-sh/nvm).
1. `cd helium-constructor`
2. `npm ci`
3. `npm run build` or `npm run watch`
4. `npm run start` to start the express server at port `:3000`

*TODO: install nodemon for development*

At this point you should be able to run the `rosetta-cli` check from above and get similiar results to the docker container. Remember, make sure to give `blockchain-node` a few minutes to warm up before it picks up blocks.

# Implementation details

### Supported currencies

- [HNT](https://www.coinbase.com/price/helium) (Helium Token)
- HST (Helium Security Token)

### Unsupported currencies
- DC (Data Credits): not implemented as they cannot be actively traded

## Data API transactions
Transactions support for reading from the Data API

| API | Implemented | TODO | Notes |
|----|:-----------:|:----:|-------|
| `payment_v1` | DONE | | |
| `payment_v2` | DONE | | |
| `reward_v1` | DONE | | |
| `reward_v2` | DONE | | |
| `security_coinbase_v1` | DONE | | |
| `security_exchange_v1` | DONE | | |
| `token_burn_v1` | | TODO | |
| `transfer_hotspot_v1` | | TODO | |
| `create_htlc_v1`* | | TODO | |
|  `redeem_htlc_v1`* | | TODO | |

### Fee-only transactions (Only recording implicit_burns for HNT deductions)

| API | Implemented | TODO | Notes |
| --- |:-----------:|:----:|-------|
| `add_gateway_v1` | DONE | | |
| `assert_location_v1` | DONE | | |
| `assert_location_v2` | DONE | | |
| `oui_v1` | | TODO | |
| `routing_v1` | | TODO | |
| `state_channel_open_v1` | | TODO | |

## Construction API transactions
Transaction support for creation via the construction API

| API | Implemented | TODO | Notes |
|-----|:-----------:|:----:|-------|
| `payment_v2` | | TODO | |
| `security_exchange_v1` | | TODO | |
| `create_htlc_v1`* | | TODO | |
| `redeem_htlc_v1`* | | TODO | |

## No Plans to Implement

| API | Implemented | Notes |
|-----|:-----------:|-------|
| `payment_v1` | NEVER | Deprecated transactions (construction API only) |
| `dc_coinbase_v1` | NEVER | DC only transaction |
| `state_channel_close_v1` | NEVER | DC only transaction |
| `gen_gateway_v1` | NEVER | Internal blockchain only |
| `poc_request_v1`| NEVER | Internal blockchain only |
| `poc_receipt_v1` | NEVER | Internal blockchain only | 
| `consensus_group_v1` | NEVER | Internal blockchain only |
| `vars_v1` | NEVER | |
| `price_oracle_v1` | NEVER | Oracle HNT value transactions | 

#### Buy me an HNT coffee :)

HNT Address: [13Ey7fZfdQB7C8FRiWfuNKshA8si7wgHVeMdGbdCkM5gyry7G88](https://explorer.helium.com/accounts/13Ey7fZfdQB7C8FRiWfuNKshA8si7wgHVeMdGbdCkM5gyry7G88)
