FROM rust:1.50

COPY . .
RUN cargo build --release --package orchestrator

CMD sh startup.sh