FROM node:15.11-alpine3.13

RUN apk update
RUN apk add --no-cache python3 make g++ curl

COPY . .
RUN npm ci
RUN chmod -R +x scripts

RUN npm run typechain

CMD sh scripts/contract-deployer.sh