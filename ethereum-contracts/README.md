# Testing contracts

## Installing dependencies

  - Truffle: `npm install -g truffle` for Mac or go to the installat
  - Ganache: Download file [here](http://truffleframework.com/ganache/)

## Instructions

You can test the contracts using Truffle Testrpc in the console or Ganache UI.

1) **Testrpc**:
    - `$ cd ethereum-contracts`
    - Delete *build* folder an type `$ truffle migrate --reset --compile-all`
    - `$ testrpc`
    - `$ truffle test` on a separate tab


2) **Ganache**:
    - `$ cd ethereum-contracts`
    - Delete *build* folder an type `$ truffle migrate --reset --compile-all`
    - Open Ganache app, open the *Settings* button and make sure it's running on port **_8545_**.
    - `$ truffle test` on a separate tab


![alt text](./img/ganache_setup.png "Ganache Setup")

![alt text](./img/settings.png "Setting ")

### Running tests of specific contracts

    `$ truffle test test/<test_name.js>`

## Debugging

Once you run the contracts you will get the `TX HASH` number of each function call. Copy it and on in a separate tab in the console:

  `$ truffle console`

  `$ debug <TX HASH>`

  ![alt text](./img/tx_hash.png "Getting TX HASH")
