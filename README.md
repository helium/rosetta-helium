[Read the wiki!](https://github.com/helium/rosetta-helium/wiki)

# Overview
Dockerized [Rosetta API](https://www.rosetta-api.org/) implementation mostly based off of [blockchain-node](https://github.com/helium/blockchain-node):
- Rosetta specs: [https://www.rosetta-api.org/](rosetta-api.org)
- This is NOT a full node, but rather works off the latest snapshot as specified in `blockchain-node`. As a result, there is currently no support for historical balances or reconciliation.
- `blockchain-node` provides the basic blockchain that the Data API reads from
- `./helium-constructor` implements a simple Express server written in TypeScript exposing [helium-js](https://github.com/helium/helium-js) for Construction API actions (transaction construction, signing mechanisms, etc)

This project was created by [@syuan100](https://github.com/syuan100) and supported, in part, by the [DeWi Grants Program](https://dewialliance.medium.com/announcing-the-inaugural-dewi-grant-recipients-56b44b9b9b66).

#### See Also

* [Rosetta Community Discussion](https://community.rosetta-api.org/t/helium-network-rosetta-gateway-implementation/529/)

# Quick setup

#### Build container from source
Mainnet:
```text
DOCKER_BUILDKIT=1 docker build . -t rosetta-helium:latest
```

Testnet:
```text
DOCKER_BUILDKIT=1 docker build . -t rosetta-helium:latest --build-arg NETWORK=testnet
```

*Note: `DOCKER_BUILDKIT=1` is not necessary but including it may reduce the image size due to the nature of the conditional build.*

#### Build container using pre-built Helium images
Mainnet:
```text
DOCKER_BUILDKIT=1 docker build . -f Dockerfile_quick -t rosetta-helium:latest
```

Testnet:
```text
DOCKER_BUILDKIT=1 docker build . -f Dockerfile_quick -t rosetta-helium:latest --build-arg NETWORK=testnet
```

#### Run container
Local data is stored in `helium-data`
```text
docker run -d --rm --init --ulimit "nofile=1000000:1000000" -v "$(pwd)/helium-data:/data" -p 8080:8080 -p 44158:44158 rosetta-helium:latest
```

It's HIGHLY recommended that you set the internal/external NAT settings through environment variables for better performance:

`NAT_INTERNAL_IP` -> `172.17.0.X` _X depends on how many docker containers you have_

`NAT_INTERNAL_PORT` -> `44158` Default port for peering

`NAT_EXTERNAL_IP` -> Your publicly accessible IP address

`NAT_EXTERNAL_PORT` -> `44158` Generally would want to keep the same port that you exposed in the command line

```text
docker run -d --rm --init --ulimit "nofile=1000000:1000000" -v "$(pwd)/helium-data:/data" -p 8080:8080 -p 44158:44158 -e NAT_INTERNAL_IP={{docker_ip}} -e NAT_INTERNAL_PORT={{docker_port}} -e NAT_EXTERNAL_IP={{public_ip}} -e NAT_EXTERNAL_PORT={{public_port}} rosetta-helium:latest
```

#### Rosetta CLI check
Mainnet:
```text
rosetta-cli check:data --configuration-file rosetta-cli-config/mainnet/config.json

rosetta-cli check:construction --configuration-file rosetta-cli-config/mainnet/config.json
```

Testnet:
```text
rosetta-cli check:data --configuration-file rosetta-cli-config/testnet/config.json

rosetta-cli check:construction --configuration-file rosetta-cli-config/testnet/config.json
```
(Please wait a few minutes for the Helium node to initialize before running this command)

[Read more on using the rosetta-cli.](https://github.com/helium/rosetta-helium/wiki/7.-Appendix:-Using-rosetta-cli-for-testing)

# Contributing
It's annoying to spin up a docker container for every change that you want to make. So for local development, it is recommended that you run each part of the implementation separately.

### rosetta-helium
1. Install [golang](https://golang.org/doc/install) if you haven't yet.
2. At the root directory, run `go run .` to start the rosetta server at port `:8080`

### blockchain-node
1. Checkout [blockchain-node](https://github.com/helium/blockchain-node/).
2. Run `make && make release` to build a release
3. Run `make start` to start blockchain-node at port `:4467`

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

| Transaction | Implemented |
|----|-----------|
| `payment_v1` | :white_check_mark: |
| `payment_v2` | :white_check_mark: |
| `reward_v1` | :white_check_mark: |
| `reward_v2` | :white_check_mark: |
| `security_coinbase_v1` | :white_check_mark: |
| `security_exchange_v1` | :white_check_mark: |
| `token_burn_v1` | :white_check_mark: |
| `transfer_hotspot_v1` | :white_check_mark: |
| `transfer_hotspot_v2` | :white_check_mark: |
|  `stake_validator_v1` | :white_check_mark: |
|  `unstake_validator_v1` | :white_check_mark: |
|  `transfer_validator_v1` | :white_check_mark: |
| `create_htlc_v1` | :white_check_mark: |
|  `redeem_htlc_v1` | :white_check_mark: |

### Fee-only transactions (Only recording implicit_burns for HNT deductions)

| Transaction | Implemented |
| --- |-----------|
| `add_gateway_v1` | :white_check_mark: |
| `assert_location_v1` | :white_check_mark: |
| `assert_location_v2` | :white_check_mark: |
| `oui_v1` | :white_check_mark: |
| `update_gateway_oui_v1` | :white_check_mark: |
| `routing_v1` | :white_check_mark: |
| `state_channel_open_v1` | :white_check_mark: |

### Pass-through transactions (No balance changes, only metadata)

| Transaction | Notes |
|-----|-------|
| `dc_coinbase_v1` | DC only transaction |
| `state_channel_close_v1` | DC only transaction |
| `gen_gateway_v1` | Internal blockchain only |
| `poc_request_v1` | Internal blockchain only |
| `poc_receipt_v1` | Internal blockchain only | 
| `consensus_group_v1` | Internal blockchain only |
| `vars_v1` | Internal blockchain only |
| `price_oracle_v1` | Oracle HNT value transactions | 

## Construction API transactions
Transaction support for creation via the construction API

| Transaction | Implemented |
|-----|-----------|
| `payment_v2` | :white_check_mark: |
| `security_exchange_v1` | :x: |
| `create_htlc_v1` | :x: |
| `redeem_htlc_v1` | :x: |
| `stake_validator_v1` | :x: |
| `unstake_validator_v1` | :x: |
| `transfer_validator_v1` | :x: |

## Additional notes
### Unstake Transaction Oddities

The `unstake_validator_v1` transaction is unique in that the balance changing portion of the transaction doesn't happen until the specified cooldown has passed. At that point, there is a callback on the ledger that records the balance change. Unfortunately, there is no way for `blockchain-node` to surface information about this balance change when inspecting a block at a particular height. This is especially important for the rosetta-cli `check:data` command to pass.

To mitigate this, we decided to store what we are calling "ghost transactions" that will be triggered when Rosetta queries a particular block. That way, Rosetta can notice the balance change and properly credit the appropriate account for reconciliation purposes.

**Example:**
- Unstake transaction dectected at block 5 with a cool down of 10 blocks (transaction fee is deducted here)
- We store a "ghost transaction" in our local BadgerDB at block 15 (current block + cool down) with the appropriate credit to the appropriate account
- When we query block 15, Rosetta will check BadgerDB for any ghost transactions and include them in the query response.