FROM fedora
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin
RUN dnf install -y git golang make gcc gcc-c++ which iproute iputils procps-ng vim tmux net-tools htop
ARG REPOFOLDER
ARG NODES
ENV NODES=$NODES
ADD $REPOFOLDER /peggy
RUN pushd /peggy/module/ && make && make install
RUN pushd /peggy/ && tests/setup-validators.sh 20
CMD pushd /peggy/ && tests/run-validators.sh 20
