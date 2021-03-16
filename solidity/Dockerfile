FROM node:15.11-alpine3.13

RUN apk add python3 make g++

COPY package.json .
COPY package-lock.json .
RUN npm ci

COPY . .
RUN chmod -R +x scripts

RUN npm run typechain

CMD sh scripts/contract-deployer.sh