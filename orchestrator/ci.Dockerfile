FROM alpine:edge

COPY target/release/orchestrator /usr/bin/orchestrator

CMD sh startup.sh