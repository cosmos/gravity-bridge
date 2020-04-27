<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

Types of tags:

"genesis": genesis state related changes
"eth-bridge-app": changes related to the application
"modules": updates to the app modules
"simulation": simulation related changes
"contracts": smart contract related changes
"docs/specs": updates to documentation and specifications
"rest": REST client changes
"cli": CLI changes

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Features

- (ethbridge) [\#75](https://github.com/cosmos/peggy/pull/75) Add burn message and functionality for burning tokens and trigerring
  events.
- (ethbridge) [\#81](https://github.com/cosmos/peggy/pull/81) Add functionality for locking up cosmos coins and detecting burned coins from ethereum
- (relayer) [\#163](https://github.com/cosmos/peggy/pull/163) Unified Relayer for chain subscription and event relay
- (ethbridge) [\#182](https://github.com/cosmos/peggy/pull/182) Remove whitelisting and prefix Eth-tokens and change API/CLI
- (relayer) [\#183](https://github.com/cosmos/peggy/pull/183) Add integrated generator for contract go bindings
- (relayer) [\#190](https://github.com/cosmos/peggy/pull/190) Support burns of native cosmos assets on Ethereum

### Improvements

- (ethbridge) [\#67](https://github.com/cosmos/peggy/pull/67) Add events to handler
- (ethbridge) [\#74](https://github.com/cosmos/peggy/pull/74) Change to use new supply module
- (testnet-contracts) [\#82](https://github.com/cosmos/peggy/pull/82) Update contracts
- (testnet-contracts) [\#89](https://github.com/cosmos/peggy/pull/89) Dynamic validator set
- (misc) [\#71](https://github.com/cosmos/peggy/pull/71) Dockerize validator and relayer for simple testing purposes
- (sdk) Update SDK version to v0.37.4
- (ethbridge) [\#100](https://github.com/cosmos/peggy/pull/100) Add Keeper to EthBridge module and use expected keepers abstraction
