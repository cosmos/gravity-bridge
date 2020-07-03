FROM fedora
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin
RUN dnf install -y git golang make gcc gcc-c++ which iproute iputils procps-ng vim tmux net-tools htop tar jq
ADD https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.9.14-6d74d1e5.tar.gz /geth/
ADD https://updates.altheamesh.com/gen_eth_key /usr/bin/
RUN chmod +x /usr/bin/gen_eth_key
RUN cd /geth && tar -xvf * && mv /geth/**/geth /usr/bin/geth
ARG REPOFOLDER
ADD $REPOFOLDER /peggy
RUN pushd /peggy/module/ && make && make install