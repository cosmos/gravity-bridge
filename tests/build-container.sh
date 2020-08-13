#!/bin/bash
set -eux

# this directy of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DOCKERFOLDER=$DIR/dockerfile
REPOFOLDER=$DIR/..

#docker system prune -a -f
# Build base container
git archive --format=tar.gz -o $DOCKERFOLDER/peggy.tar.gz --prefix=peggy/ HEAD
pushd $DOCKERFOLDER
docker build -t peggy-base .
