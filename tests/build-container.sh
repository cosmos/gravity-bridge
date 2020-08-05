#!/bin/bash
set -eux

#docker system prune -a -f
# Build base container
docker build -f ./tests/dockerfile/base.Dockerfile -t peggy-base .
