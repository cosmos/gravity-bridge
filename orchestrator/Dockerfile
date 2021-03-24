FROM rust:1.50

# haven't found sccache to improve docker build times
#RUN cargo install sccache --verbose
#ENV RUSTC_WRAPPER="/usr/local/cargo/bin/sccache"

COPY . .
RUN cargo build --workspace --release

#ENV COSMOS_KEY="UNSET"
#ENV ETH_PRIVATE_KEY="UNSET"
#ENV ETH_PRIVATE_KEY="UNSET"
#ENV COSMOS_RPC="UNSET"
#ENV COSMOS_GRPC="UNSET"
#ENV ETH_RPC="UNSET"
#ENV DENOM="UNSET"
#ENV CONTRACT_ADDR="UNSET"

CMD sh startup.sh