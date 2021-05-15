First run

1. sudo dnf install openssl-devel perl
2. cargo check --all
3. cargo test --all

Regenerate proto after updated proto files

```
cd proto_build
cargo run
```
