FROM peggy-base
ARG REPOFOLDER
COPY . /peggy
RUN pushd /peggy/module/ && make && make install
RUN pushd /peggy/ && tests/setup-validators.sh 3
CMD pushd /peggy/ && tests/run-testnet.sh 3