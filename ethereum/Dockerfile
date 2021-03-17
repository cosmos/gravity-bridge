FROM ethereum/client-go:stable

COPY tests/assets/ETHGenesis.json ETHGenesis.json

RUN geth --identity "PeggyTestnet" \
    --nodiscover \
    --networkid 15 init ETHGenesis.json

ENTRYPOINT geth --identity "PeggyTestnet" --nodiscover \
           --networkid 15 \
           --mine \
           --rpc \
           --http \
           --miner.threads=1 \
           --nousb \
           --verbosity=5 \
           --miner.etherbase=0xBf660843528035a5A4921534E156a27e64B231fE