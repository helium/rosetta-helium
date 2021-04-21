# Compile golang 
FROM ubuntu:20.04 as golang-builder

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app
WORKDIR /app

RUN apt-get update && apt-get install -y curl make gcc g++ git
ENV GOLANG_VERSION 1.15.5
ENV GOLANG_DOWNLOAD_SHA256 9a58494e8da722c3aef248c9227b0e9c528c7318309827780f16220998180a0d
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
  && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
  && tar -C /usr/local -xzf golang.tar.gz \
  && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Compile erlang 
FROM ubuntu:20.04 as helium-builder

ENV DEBIAN_FRONTEND noninteractive
ENV PATH /root/.cargo/bin:${PATH}

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app
WORKDIR /app

RUN apt-get update && apt-get install -y curl make gcc g++ git gnupg2 wget
RUN wget https://packages.erlang-solutions.com/erlang-solutions_2.0_all.deb && dpkg -i erlang-solutions_2.0_all.deb
RUN apt-get update \
    && apt-get install -y esl-erlang=1:22.3.4.1-1 cmake libsodium-dev libssl-dev build-essential
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

# Compile rosetta-ethereum
FROM golang-builder as rosetta-builder

# Use native remote build context to build in any directory
COPY . rosetta-helium-builder
RUN cd rosetta-helium-builder \
  && go build -o rosetta-helium

RUN mv rosetta-helium-builder/rosetta-helium /app/rosetta-helium \
  && rm -rf rosetta-helium-builder

FROM helium-builder as blockchain-node-builder

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app \
  && mkdir -p /data \
  && chown -R nobody:nogroup /data

# Copy Makefile
COPY ./docker/Makefile /app

# VERSION: blockchain-node 
RUN git clone https://github.com/syuan100/blockchain-node \
   && cd blockchain-node \
   && git checkout ced91e1a3ef1d1942022c4585fe5d71e1117ea41

RUN cd blockchain-node \
  && make && make release PROFILE=devib

RUN cat ./blockchain-node/config/sys.config | grep -oP '(?<=\{blessed_snapshot_block_height\, ).*?(?=\})' > lbs.txt

COPY --from=rosetta-builder /app/rosetta-helium /app/rosetta-helium

RUN chmod -R 755 /app/*

WORKDIR /app

CMD ["make", "start", "PROFILE=devib"]
