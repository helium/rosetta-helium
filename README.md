**THIS IS NOT PRODUCTION READY. USE AT YOUR OWN RISK.**

## Overview
Bare bones Rosetta API implementation of the Helium `blockchain-node`

#### Build container
```text
docker build . -t rosetta-helium:latest
```

#### Run container
Local data is stored in `helium-data`
```text
docker run -d --rm --ulimit "nofile=100000:100000" -v "$(pwd)/helium-data:/app/blockchain-node/_build/dev/rel/blockchain_node/data" -p 8080:8080 -p 44158:44158 rosetta-helium:latest
```

#### Rosetta CLI check
```text
rosetta-cli check:data --configuration-file rosetta-cli-config/mainnet/config.json
```
(Please wait a few minutes for the Helium node to initialize before running this command)

#### Buy me an HNT coffee :)
HNT Address: 13Ey7fZfdQB7C8FRiWfuNKshA8si7wgHVeMdGbdCkM5gyry7G88
