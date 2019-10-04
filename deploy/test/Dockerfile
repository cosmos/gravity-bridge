FROM golang:1.13

WORKDIR /go/src/github.com/cosmos/peggy

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make install

RUN ebd init local --chain-id=peggy
RUN echo password | ebcli keys import validator ./deploy/test/test_validator_key
RUN ebd add-genesis-account $(ebcli keys show validator -a) 1000000000stake,1000000000atom
RUN echo password | ebd gentx --name validator
RUN ebd collect-gentxs