FROM rust:1.50

COPY . .
RUN cargo build --bin orchestrator

CMD sh startup.sh