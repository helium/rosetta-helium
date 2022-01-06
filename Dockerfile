ARG NETWORK=mainnet

FROM ubuntu:20.04 as node-builder

ARG BUILD_VERSION="23.3.4.7"
ARG BUILD_SHA256="37e39a43c495861ce69de06e1a013a7eac81d15dc6eebd2d2022fd68791f4b2d"
ENV OTP_VERSION=$BUILD_VERSION \
    REBAR3_VERSION="3.16.1"

LABEL org.opencontainers.image.version=$OTP_VERSION

# We'll install the build dependencies, and purge them on the last step to make
# sure our final image contains only what we've just built:
RUN set -xe \
	&& OTP_DOWNLOAD_URL="https://github.com/erlang/otp/releases/download/OTP-${BUILD_VERSION}/otp_src_${BUILD_VERSION}.tar.gz" \
	&& OTP_DOWNLOAD_SHA256="${BUILD_SHA256}" \
	&& fetchDeps=' \
		curl \
		ca-certificates' \
	&& apt-get update \
	&& apt-get install -y --no-install-recommends $fetchDeps \
	&& curl -fSL -o otp-src.tar.gz "$OTP_DOWNLOAD_URL" \
	&& echo "$OTP_DOWNLOAD_SHA256  otp-src.tar.gz" | sha256sum -c - \
	&& runtimeDeps=' \
		libssl1.1 \
	' \
	&& buildDeps=' \
		autoconf \
		dpkg-dev \
		gcc \
		g++ \
		make \
		libncurses-dev \
		libssl-dev \
	' \
	&& apt-get install -y --no-install-recommends $runtimeDeps \
	&& apt-get install -y --no-install-recommends $buildDeps \
	&& export ERL_TOP="/usr/src/otp_src_${OTP_VERSION%%@*}" \
	&& mkdir -vp $ERL_TOP \
	&& tar -xzf otp-src.tar.gz -C $ERL_TOP --strip-components=1 \
	&& rm otp-src.tar.gz \
	&& ( cd $ERL_TOP \
	  && ./otp_build autoconf \
	  && gnuArch="$(dpkg-architecture --query DEB_HOST_GNU_TYPE)" \
	  && ./configure --build="$gnuArch" \
	  && make -j$(nproc) \
	  && make install ) \
	&& find /usr/local -name examples | xargs rm -rf \
	&& REBAR3_DOWNLOAD_URL="https://github.com/erlang/rebar3/archive/${REBAR3_VERSION}.tar.gz" \
	&& REBAR3_DOWNLOAD_SHA256="a14711b09f6e1fc1b080b79d78c304afebcbb7fafed9d0972eb739f0ed89121b" \
	&& mkdir -p /usr/src/rebar3-src \
	&& curl -fSL -o rebar3-src.tar.gz "$REBAR3_DOWNLOAD_URL" \
	&& echo "$REBAR3_DOWNLOAD_SHA256 rebar3-src.tar.gz" | sha256sum -c - \
	&& tar -xzf rebar3-src.tar.gz -C /usr/src/rebar3-src --strip-components=1 \
	&& rm rebar3-src.tar.gz \
	&& cd /usr/src/rebar3-src \
	&& HOME=$PWD ./bootstrap \
	&& install -v ./rebar3 /usr/local/bin/ \
	&& rm -rf /usr/src/rebar3-src \
	&& apt-get purge -y --auto-remove $buildDeps $fetchDeps \
	&& rm -rf $ERL_TOP /var/lib/apt/lists/*

RUN set -xe \
	    && apt update \
	    && apt-get install -y --no-install-recommends \
     	        libssl-dev make automake autoconf libncurses5-dev gcc \
	        libdbus-1-dev libbz2-dev bison flex libgmp-dev liblz4-dev \
	        libsodium-dev sed wget curl build-essential libtool git \
		ca-certificates \
	    && mkdir -p /opt/cmake \
	    && wget -O /opt/cmake/cmake.tgz \
	        https://github.com/Kitware/CMake/releases/download/v3.21.3/cmake-3.21.3-linux-x86_64.tar.gz \
            && tar -zxf /opt/cmake/cmake.tgz -C /opt/cmake --strip-components=1

# Install Rust toolchain
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

ENV CC=gcc CXX=g++ CFLAGS="-U__sun__" \
    ERLANG_ROCKSDB_OPTS="-DWITH_BUNDLE_SNAPPY=ON -DWITH_BUNDLE_LZ4=ON" \
    ERL_COMPILER_OPTIONS="[deterministic]" \
    PATH="/root/.cargo/bin:/opt/cmake/bin:$PATH" \
    RUSTFLAGS="-C target-feature=-crt-static"

WORKDIR /usr/src/

# Add our code
RUN git clone https://github.com/helium/blockchain-node.git


FROM node-builder AS node-mainnet

ARG BUILD_TARGET=docker_rosetta

WORKDIR /usr/src/blockchain-node

RUN ./rebar3 as ${BUILD_TARGET} tar -n blockchain_node

RUN mkdir -p /opt/blockchain_node \
	&& tar -zxvf _build/${BUILD_TARGET}/rel/*/*.tar.gz -C /opt/blockchain_node


FROM node-builder AS node-testnet

ARG BUILD_TARGET=docker_rosetta_testnet

WORKDIR /usr/src/blockchain-node

RUN ./rebar3 as ${BUILD_TARGET} tar -n blockchain_node

RUN mkdir -p /opt/blockchain_node \
	&& tar -zxvf _build/${BUILD_TARGET}/rel/*/*.tar.gz -C /opt/blockchain_node


FROM ubuntu:20.04 AS rosetta-builder

RUN set -xe \
	&& ulimit -n 100000 \
        && apt update \
 	&& apt install -y --no-install-recommends libdbus-1-3 libgmp10 libsodium23 \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /opt/blockchain_node

ENV COOKIE=blockchain_node \
    # Write files generated during startup to /tmp
    RELX_OUT_FILE_PATH=/tmp \
    # add to path, for easy exec interaction
    PATH=/sbin:/bin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH:/opt/blockchain_node/bin

WORKDIR /src

RUN apt update \
      && apt install -y --no-install-recommends \
         curl ca-certificates git \      
      && curl -L https://golang.org/dl/go1.17.1.linux-amd64.tar.gz | tar xzf -

ENV PATH="/src/go/bin:$PATH" \
    CGO_ENABLED=0

COPY . rosetta-helium

RUN cd rosetta-helium && go build -o rosetta-helium


FROM node-${NETWORK} as rosetta-helium-final

ARG NETWORK
ARG DEBIAN_FRONTEND=noninteractive

EXPOSE 8080
EXPOSE 44158

RUN apt update \
    && apt install -y --no-install-recommends \
         ca-certificates git npm

WORKDIR /app

COPY --from=rosetta-builder /src/rosetta-helium/rosetta-helium rosetta-helium
COPY --from=rosetta-builder /src/rosetta-helium/ghost-transactions ghost-transactions
COPY --from=rosetta-builder /src/rosetta-helium/docker/${NETWORK}.sh start.sh
COPY --from=rosetta-builder /src/rosetta-helium/helium-constructor helium-constructor

RUN cd helium-constructor \
      && npm install \
      && npm run build \
      && chmod +x /app/start.sh

ENTRYPOINT ["/app/start.sh"]
