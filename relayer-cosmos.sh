#!/bin/bash
until ebrelayer init cosmos tcp://localhost:26657 https://dai.poa.network 0x44b30c29e64031A021AB6c8DE475AA5D9AC22740; do
    echo "Server 'cosmos relayer' crashed with exit code $?.  Respawning.." >&2
    sleep 1
done
