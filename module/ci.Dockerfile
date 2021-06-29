FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates

EXPOSE 26656 26657 1317 9090

COPY build/gravity /usr/bin/gravity

CMD ["gravity", "start"]