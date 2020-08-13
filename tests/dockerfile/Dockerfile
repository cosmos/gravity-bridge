FROM fedora
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin
RUN dnf install -y git golang make gcc gcc-c++ which iproute iputils procps-ng vim tmux net-tools htop tar jq npm openssl-devel
ADD https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.9.14-6d74d1e5.tar.gz /geth/
ADD https://updates.altheamesh.com/gen_eth_key /usr/bin/
RUN chmod +x /usr/bin/gen_eth_key
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
RUN cd /geth && tar -xvf * && mv /geth/**/geth /usr/bin/geth
ADD peggy.tar.gz /
RUN pushd /peggy/tests/test-runner && PATH=$PATH:$HOME/.cargo/bin cargo build
RUN pushd /peggy/module/ && make && make install
RUN pushd /peggy/solidity/ && npm install