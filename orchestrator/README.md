# Orchestrator

The Orchestrator is responsible for both relaying transactions to and from evm based chains chains.

## Developer Docs

### Compile Protobuf

To compile new protobuf definitions follow the below steps:

```sh
cd proto_build
cargo run
```

After you have run the above two commands, go to `gravity_proto/prost`. You will need to delete all the files in prost directory, other than `gravity.v1.rs`.
