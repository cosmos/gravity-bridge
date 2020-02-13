#!/bin/bash
until ebrelayer init cosmos tcp://167.99.167.78:26657 https://dai.poa.network 0x4484aaD19922304C4f3A6aA1D0D65C79266e0d11; do
    echo "Server 'cosmos relayer' crashed with exit code $?.  Respawning.." >&2
    sleep 1
done
