FROM ubuntu:20.04 as node-builder
RUN apt-get update && apt-get install -y curl make gcc g++ git gnupg2 wget

ENV DEBIAN_FRONTEND noninteractive
ENV PATH /root/.cargo/bin:${PATH}

RUN wget https://packages.erlang-solutions.com/erlang-solutions_2.0_all.deb && dpkg -i erlang-solutions_2.0_all.deb
RUN apt-get update \
    && apt-get install -y esl-erlang=1:22.3.4.1-1 cmake libsodium-dev libssl-dev build-essential
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

WORKDIR /usr/src

# Add our code
RUN git clone https://github.com/syuan100/blockchain-node \
   && cd blockchain-node \
   && git checkout 52586988fb4ed4ce81333bf5d4cedcf17fa86292

WORKDIR /usr/src/blockchain-node

RUN ./rebar3 as devib tar
RUN mkdir -p /opt/blockchain-node-build
RUN tar -zxvf _build/devib/rel/*/*.tar.gz -C /opt/blockchain-node-build

FROM ubuntu:20.04 as rosetta-builder
RUN mkdir -p /app \
  && chown -R nobody:nogroup /app
WORKDIR /app

RUN apt-get update && apt-get install -y curl make gcc g++ git gnupg2 wget
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

COPY . rosetta-helium-builder

WORKDIR /app/rosetta-helium-builder

RUN go build -o rosetta-helium
RUN mv rosetta-helium /app/rosetta-helium

WORKDIR /app

RUN rm -rf rosetta-helium-builder

FROM ubuntu:20.04 as runner
RUN apt-get update && apt-get install -y gnupg2 wget

ENV DEBIAN_FRONTEND noninteractive

RUN wget https://packages.erlang-solutions.com/erlang-solutions_2.0_all.deb && dpkg -i erlang-solutions_2.0_all.deb
RUN apt-get update \
    && apt-get install -y esl-erlang=1:22.3.4.1-1

RUN ulimit -n 100000

ENV COOKIE=blockchain-node \
    # Write files generated during startup to /tmp
    RELX_OUT_FILE_PATH=/tmp \
    # add miner to path, for easy interactions
    PATH=$PATH:/app/blockchain-node/bin

COPY --from=node-builder /opt/blockchain-node-build /app/blockchain-node
COPY --from=rosetta-builder /app/rosetta-helium /app/rosetta-helium

RUN cat /app/blockchain-node/config/sys.config | grep -oP '(?<=\{blessed_snapshot_block_height\, ).*?(?=\})' > lbs.txt

COPY ./docker/start.sh /app

RUN chmod +x /app/start.sh

CMD ["/app/start.sh"]