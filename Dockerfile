FROM erlang:22.3.2-alpine as node-builder
RUN apk add --no-cache --update \
    git tar build-base linux-headers autoconf automake libtool pkgconfig \
    dbus-dev bzip2 bison flex gmp-dev cmake lz4 libsodium-dev openssl-dev \
    sed wget curl

ENV CC=gcc CXX=g++ CFLAGS="-U__sun__" \
    ERLANG_ROCKSDB_OPTS="-DWITH_BUNDLE_SNAPPY=ON -DWITH_BUNDLE_LZ4=ON" \
    ERL_COMPILER_OPTIONS="[deterministic]" \
    PATH="/root/.cargo/bin:$PATH" \
    RUSTFLAGS="-C target-feature=-crt-static"

# Install Rust toolchain
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

WORKDIR /usr/src

# Add our code
RUN git clone https://github.com/syuan100/blockchain-node \
   && cd blockchain-node \
   && git checkout 52586988fb4ed4ce81333bf5d4cedcf17fa86292

WORKDIR /usr/src/blockchain-node

RUN ./rebar3 as devib tar
RUN mkdir -p /opt/blockchain-node-build
RUN tar -zxvf _build/devib/rel/*/*.tar.gz -C /opt/blockchain-node-build

FROM erlang:22.3.2-alpine as rosetta-builder
RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go
ENV PATH="/usr/local/go/bin:$PATH" \
    GOPATH=/opt/go/ \
    PATH=$PATH:$GOPATH/bin 

WORKDIR /app

COPY . rosetta-helium-builder

WORKDIR /app/rosetta-helium-builder

RUN go build -o rosetta-helium
RUN mv rosetta-helium /app/rosetta-helium

WORKDIR /app

RUN rm -rf rosetta-helium-builder

FROM erlang:22.3.2-alpine as runner
RUN apk add --no-cache --update --virtual .build-deps bash gcc musl-dev openssl go grep ncurses dbus gmp libsodium
ENV PATH="/usr/local/go/bin:$PATH" \
    GOPATH=/opt/go/ \
    PATH=$PATH:$GOPATH/bin 

RUN ulimit -n 100000

ENV COOKIE=blockchain-node \
    # Write files generated during startup to /tmp
    RELX_OUT_FILE_PATH=/tmp \
    # add miner to path, for easy interactions
    PATH=$PATH:/app/blockchain-node/bin \
    CGO_ENABLED=0

COPY --from=node-builder /opt/blockchain-node-build /app/blockchain-node
COPY --from=rosetta-builder /app/rosetta-helium /app/rosetta-helium

RUN cat /app/blockchain-node/config/sys.config | grep -oP '(?<=\{blessed_snapshot_block_height\, ).*?(?=\})' > lbs.txt

COPY ./docker/start.sh /app

RUN chmod +x /app/start.sh

CMD ["/app/start.sh"]