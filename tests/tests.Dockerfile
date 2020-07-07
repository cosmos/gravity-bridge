FROM peggy-base
COPY . /peggy
ENV NODES=3
CMD pushd /peggy/ && tests/reload-code.sh $NODES
# RUN pushd /peggy/module/ && make && make install
# RUN pushd /peggy/ && tests/setup-validators.sh $NODES
# CMD pushd /peggy/ && tests/run-testnet.sh $NODES