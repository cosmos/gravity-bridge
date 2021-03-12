FROM fedora:33 AS go

ENV GOPATH=/go
ENV PATH=$PATH:/go/bin:/usr/local/go/bin
ENV GOPROXY=https://proxy.golang.org
ENV PATH=$PATH

RUN dnf install -y git make gcc gcc-c++ which iproute iputils procps-ng vim-minimal tmux net-tools htop tar jq npm openssl-devel perl
ADD https://golang.org/dl/go1.15.6.linux-amd64.tar.gz /go/
RUN cd /go && tar -xvf * && mv /go/**/ /usr/local/

FROM go AS module

COPY . .
CMD 'make && make install'

