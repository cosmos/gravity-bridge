FROM ethereum/client-go:stable AS base-env

RUN apk add git make npm python
RUN npm install -g ts-node typescript
#ADD https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.9.14-6d74d1e5.tar.gz /geth/
#RUN cd /geth && tar -xvf * && mv /geth/**/geth /usr/bin/geth

FROM base-env AS ethereum

COPY . .
RUN npm ci
RUN npm typechain

