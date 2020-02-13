 #!/bin/bash
 until ebrelayer init ethereum  ws://xdai-rpc-parity-005.poa.network:8546 0x7cA01596f991e464C2DD4E9547Bc152291176D71 validator --chain-id=peggy; do
    echo "Server 'ethereum relayer' crashed with exit code $?.  Respawning.." >&2
    sleep 1
done