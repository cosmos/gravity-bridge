FROM fedora
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin
RUN dnf install -y git golang make gcc gcc-c++ which iproute iputils procps-ng vim tmux net-tools htop tar jq npm openssl-devel
# only required for deployment script
RUN npm install -g ts-node && npm install -g typescript
ADD https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.9.14-6d74d1e5.tar.gz /geth/
RUN cd /geth && tar -xvf * && mv /geth/**/geth /usr/bin/geth
# super simple rust program that dumps out an eth private key
ADD https://updates.altheamesh.com/gen_eth_key /usr/bin/
RUN chmod +x /usr/bin/gen_eth_key
# rust for the testing suite and the relayers / validator daemons
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
# the actual source code for this repo, this **only** includes checked in files!
# this is a bit of a pain but it does speed things up a lot
ADD peggy.tar.gz /
# build steps for all codebases in this repo, must be below the add statement
RUN pushd /peggy/tests/test-runner && PATH=$PATH:$HOME/.cargo/bin cargo build
RUN pushd /peggy/module/ && make && make install
RUN pushd /peggy/solidity/ && npm install