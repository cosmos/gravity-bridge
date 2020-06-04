FROM fedora
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin
RUN dnf install -y git golang make gcc gcc-c++ which iproute iputils procps-ng vim tmux net-tools htop tar
ADD https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.9.14-6d74d1e5.tar.gz /geth/
RUN cd /geth && tar -xvf * && mv /geth/**/geth /usr/bin/geth
ARG REPOFOLDER
ARG NODES
ENV NODES=$NODES
ADD $REPOFOLDER /peggy
RUN pushd /peggy/module/ && make && make install
RUN pushd /peggy/ && tests/setup-validators.sh 3
CMD pushd /peggy/ && tests/run-testnet.sh 3
