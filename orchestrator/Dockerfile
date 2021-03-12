FROM fedora:33 AS rust

ENV PATH=$PATH:$HOME/.cargo/bin

RUN dnf install -y git make gcc gcc-c++ which iproute iputils procps-ng vim-minimal tmux net-tools htop tar jq npm openssl-devel perl
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

FROM rust AS orchestrator

COPY . .
CMD 'cargo build --all --release'
