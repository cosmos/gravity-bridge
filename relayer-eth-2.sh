 #!/bin/bash
 until ebrelayer init ethereum  ws://xdai-rpc-parity-005.poa.network:8546 0x44b30c29e64031A021AB6c8DE475AA5D9AC22740 validator --chain-id=peggy; do
    echo "Server 'ethereum relayer' crashed with exit code $?.  Respawning.." >&2
    sleep 1
done