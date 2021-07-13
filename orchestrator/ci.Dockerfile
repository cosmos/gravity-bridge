FROM alpine:edge

COPY target/release/orchestrator /usr/bin/orchestrator
COPY startup.sh startup.sh

CMD sh startup.sh