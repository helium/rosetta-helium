**THIS IS NOT PRODUCTION READY. USE AT YOUR OWN RISK.**

[Read the wiki!](https://github.com/syuan100/rosetta-helium/wiki)

# Overview
Dockerized [Rosetta API](https://www.rosetta-api.org/) implementation mostly based off of [blockchain-node](https://github.com/helium/blockchain-node):
- Rosetta specs: [https://www.rosetta-api.org/](rosetta-api.org)
- This is NOT a full node, but rather works off the latest snapshot as specified in `blockchain-node`. As a result, there is currently no support for historical balances or reconciliation.
- `blockchain-node` provides the basic blockchain that the Data API reads from
- `./helium-constructor` implements a simple Express server written in TypeScript exposing [helium-js](https://github.com/helium/helium-js) for Construction API actions (transaction construction, signing mechanisms, etc)

## Support This Project

#### Buy me an HNT coffee :)

HNT Address: [13Ey7fZfdQB7C8FRiWfuNKshA8si7wgHVeMdGbdCkM5gyry7G88](https://explorer.helium.com/accounts/13Ey7fZfdQB7C8FRiWfuNKshA8si7wgHVeMdGbdCkM5gyry7G88)

# Quick setup

#### Build container
```text
docker build . -t rosetta-helium:latest
```

#### Run container
Local data is stored in `helium-data`
```text
docker run -d --rm --ulimit "nofile=1000000:1000000" -v "$(pwd)/helium-data:/data" -p 8080:8080 -p 44158:44158 rosetta-helium:latest
```

#### Rosetta CLI check
```text
rosetta-cli check:data --configuration-file rosetta-cli-config/mainnet/config.json

rosetta-cli check:construction --configuration-file rosetta-cli-config/mainnet/config.json
```
(Please wait a few minutes for the Helium node to initialize before running this command)

[Read more on using the rosetta-cli.](https://github.com/syuan100/rosetta-helium/wiki/7.-Appendix:-Using-rosetta-cli-for-testing)

# Contributing
It's annoying to spin up a docker container for every change that you want to make. So for local development, it is recommended that you run each part of the implementation separately.

### rosetta-helium
1. Install [golang](https://golang.org/doc/install) if you haven't yet.
2. At the root directory, run `go run .` to start the rosetta server at port `:8080`

### blockchain-node
1. Checkout my [custom version of blockchain-node](https://github.com/syuan100/blockchain-node/tree/syuan100-rosetta-api) that accounts for [implicit burn](https://docs.helium.com/blockchain/transaction-fees/) events.
2. Run `make && make release PROFILE=devib` to build a release
3. Run `make start PROFILE=devib` to start blockchain-node at port `:4467`

### helium-constructor
1. Install `node`. I prefer [nvm](https://github.com/nvm-sh/nvm).
1. `cd helium-constructor`
2. `npm ci`
3. `npm run build` or `npm run watch`
4. `npm run nodemon` to start the express server at port `:3000`

At this point you should be able to run the `rosetta-cli` check from above and get similiar results to the docker container. Remember, make sure to give `blockchain-node` a few minutes to warm up before it picks up blocks.

# Implementation details

### Supported currencies

- [HNT](https://www.coinbase.com/price/helium) (Helium Token)
- HST (Helium Security Token)

### Unsupported currencies
- DC (Data Credits): not implemented as they cannot be actively traded

## Data API transactions
Transactions support for reading from the Data API

| API | Implemented |
|----|-----------|
| `payment_v1` | :white_check_mark: |
| `payment_v2` | :white_check_mark: |
| `reward_v1` | :white_check_mark: |
| `reward_v2` | :white_check_mark: |
| `security_coinbase_v1` | :white_check_mark: |
| `security_exchange_v1` | :white_check_mark: |
| `token_burn_v1` | :white_check_mark: |
| `transfer_hotspot_v1` | :white_check_mark: |
|  `stake_validator_v1` | :white_check_mark: |
|  `unstake_validator_v1` | :white_check_mark: |
|  `transfer_validator_v1` | :white_check_mark: |
| `create_htlc_v1`* | Considered |
|  `redeem_htlc_v1`* | Considered |

### Fee-only transactions (Only recording implicit_burns for HNT deductions)

| API | Implemented |
| --- |-----------|
| `add_gateway_v1` | :white_check_mark: |
| `assert_location_v1` | :white_check_mark: |
| `assert_location_v2` | :white_check_mark: |
| `oui_v1` | :white_check_mark: |
| `routing_v1` | :white_check_mark: |
| `state_channel_open_v1` | :white_check_mark: |

## Construction API transactions
Transaction support for creation via the construction API

| API | Implemented |
|-----|-----------|
| `payment_v2` | :white_check_mark: |
| `security_exchange_v1` | CONSIDERED |
| `create_htlc_v1`* | CONSIDERED |
| `redeem_htlc_v1`* | CONSIDERED |
| `stake_validator_v1`* | CONSIDERED |
| `unstake_validator_v1`* | CONSIDERED |
| `transfer_validator_v1`* | CONSIDERED |

## No Plans to Implement

| API | Implemented | Notes |
|-----|:-----------:|-------|
| `payment_v1` | NEVER | Deprecated transaction (read-only reference in Data API) |
| `dc_coinbase_v1` | NEVER | DC only transaction |
| `state_channel_close_v1` | NEVER | DC only transaction |
| `gen_gateway_v1` | NEVER | Internal blockchain only |
| `poc_request_v1`| NEVER | Internal blockchain only |
| `poc_receipt_v1` | NEVER | Internal blockchain only | 
| `consensus_group_v1` | NEVER | Internal blockchain only |
| `vars_v1` | NEVER | Internal blockchain only |
| `price_oracle_v1` | NEVER | Oracle HNT value transactions | 
