FROM ubuntu:bionic as node-builder

ENV DEBIAN_FRONTEND noninteractive
ENV PATH /root/.cargo/bin:${PATH}
ENV CC=gcc CXX=g++ CFLAGS="-U__sun__" \
    ERLANG_ROCKSDB_OPTS="-DWITH_BUNDLE_SNAPPY=ON -DWITH_BUNDLE_LZ4=ON" \
    ERL_COMPILER_OPTIONS="[deterministic]" \
    PATH="/root/.cargo/bin:$PATH" \
    RUSTFLAGS="-C target-feature=-crt-static"

# install erlang
RUN apt-get update && apt-get install -y build-essential cmake autoconf automake bison flex libtool libsodium-dev libssl-dev libdbus-1-dev libgmp-dev libbz2-dev liblz4-dev
RUN wget https://packages.erlang-solutions.com/erlang-solutions_2.0_all.deb && dpkg -i erlang-solutions_2.0_all.deb
RUN apt-get update && apt-get install -y esl-erlang=1:22.3.4.1-1

# install rust toolchain
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

WORKDIR /usr/src

# Add our code
RUN git clone https://github.com/helium/blockchain-node \
   && cd blockchain-node \
   && git checkout master

WORKDIR /usr/src/blockchain-node

RUN ./rebar3 as docker_rosetta tar
RUN mkdir -p /opt/blockchain-node-build \
 && tar -zxvf _build/docker_rosetta/rel/*/*.tar.gz -C /opt/blockchain-node-build

####
FROM ubuntu:bionic as rosetta-builder
RUN apt-get update && apt-get install -y curl make gcc g++ git gnupg2 wget software-properties-common \
    && add-apt-repository ppa:longsleep/golang-backports && apt-get install -y golang-go
ENV PATH="/usr/local/go/bin:$PATH" \
    GOPATH=/opt/go/ \
    PATH=$PATH:$GOPATH/bin 

WORKDIR /app/builder
RUN git clone https://github.com/helium/rosetta-helium \
    && cd rosetta-helium \
    && git checkout main \
    && go build -o rosetta-helium
RUN cd rosetta-helium \
    && mv rosetta-helium /app \
    && mv docker/start.sh /app \
    && cp -R helium-constructor /app

RUN rm -rf /app/builder/rosetta-helium

####
FROM ubuntu:bionic as runner
RUN apt-get update && apt-get install -y libsodium-dev gcc
RUN apt-get install -y openssl grep dbus libgmp3-dev npm

ENV COOKIE=blockchain-node \
    # Write files generated during startup to /tmp
    RELX_OUT_FILE_PATH=/tmp \
    # add miner to path, for easy interactions
    PATH=$PATH:/app/blockchain-node/bin \
    CGO_ENABLED=0

COPY --from=node-builder /opt/blockchain-node-build /app/blockchain-node
COPY --from=rosetta-builder /app/rosetta-helium /app/rosetta-helium
COPY --from=rosetta-builder /app/start.sh /app/start.sh
COPY --from=rosetta-builder /app/helium-constructor /app/helium-constructor

RUN cd /app/helium-constructor \
    && npm install \
    && npm run build

RUN chmod +x /app/start.sh \
 && cat /app/blockchain-node/config/sys.config | grep -oP '(?<=\{blessed_snapshot_block_height\, ).*?(?=\})' > lbs.txt

CMD ["/app/start.sh"]
