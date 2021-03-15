FROM rust:1.50

# haven't found sccache to improve docker build times
#RUN cargo install sccache --verbose
#ENV RUSTC_WRAPPER="/usr/local/cargo/bin/sccache"

COPY . .
RUN cargo build --workspace --release

CMD ["cargo run --release test_runner"]