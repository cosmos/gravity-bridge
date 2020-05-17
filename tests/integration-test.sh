#!/bin/bash
set -eux
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DOCKERFOLDER=$DIR
REPOFOLDER=$DIR/..
NODES=50
# clean up docker images system wide, this does bust caching but it also
# keeps storage requirements reasonable. Without it you won't be able to
# run the test again and again without running out of root disk space
#docker system prune -a -f

pushd $REPOFOLDER
time docker build -f $DOCKERFOLDER/Dockerfile -t peggy-test .
time docker run --cap-add=NET_ADMIN -it peggy-test
popd