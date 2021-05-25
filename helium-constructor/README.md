# Helium Transaction Constructor
Node.js endpoint that uses [helium-js](https://github.com/helium/helium-js/) to construct and sign transactions.

## Usage
`helium-constructor` will normally be referenced via the Dockerfile and used as part of the docker-compose.yml in the top level repo. Rosetta endpoints that will use this middleware:

- /construction/combine
- /construction/hash
- /construction/parse
- /construction/payloads
- /construction/submit


## Development setup
Install dependencies
```
npm ci
```

Build
```
npm run build
```

Start server
```
npm run start
```