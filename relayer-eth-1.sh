 #!/bin/bash
 until ebrelayer init ethereum  ws://localhost:8546 0x4484aaD19922304C4f3A6aA1D0D65C79266e0d11 validator --keyring-backend test --make-claims=true --chain-id=peggy; do
    echo "Server 'ethereum relayer with claims' crashed with exit code $?.  Respawning.." >&2
    sleep 1
done