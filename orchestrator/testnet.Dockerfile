# Reference: https://www.lpalmieri.com/posts/fast-rust-docker-builds/

FROM rust:1.52 as cargo-chef-rust
RUN apt-get install bash
RUN cargo install cargo-chef

FROM cargo-chef-rust as planner
WORKDIR app
# We only pay the installation cost once,
# it will be cached from the second build onwards
# To ensure a reproducible build consider pinning
# the cargo-chef version with `--version X.X.X`
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

FROM cargo-chef-rust as cacher
WORKDIR app
COPY --from=planner /app/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

FROM cargo-chef-rust as builder
WORKDIR app
COPY . .
# Copy over the cached dependencies
COPY --from=cacher /app/target target
COPY --from=cacher /usr/local/cargo /usr/local/cargo
RUN cargo build --manifest-path=test_runner/Cargo.toml --release --bin test_runner

FROM cargo-chef-rust as runtime
WORKDIR app
COPY test_runner/startup.sh startup.sh
COPY --from=builder /app/target/release/test_runner /usr/local/bin
CMD sh startup.sh