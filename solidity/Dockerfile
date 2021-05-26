FROM node:15.11-alpine3.13

RUN apk update
RUN apk add --no-cache python3 make g++ curl

COPY package.json package.json
COPY package-lock.json package-lock.json
RUN npm ci

COPY . .
RUN npm ci
RUN chmod -R +x scripts

RUN npm run typechain

CMD npx ts-node \
    contract-deployer.ts \
    --cosmos-node="http://gravity0:26657" \
    --eth-node="http://ethereum:8545" \
    --eth-privkey="0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7" \
    --contract=artifacts/contracts/Gravity.sol/Gravity.json \
    --test-mode=true