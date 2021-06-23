#!/bin/bash
set -eux
# the directory of this script, useful for allowing this script
# to be run with any PWD
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR
bash all-up-test.sh
bash all-up-test.sh VALSET_STRESS
bash all-up-test.sh BATCH_STRESS
bash all-up-test.sh HAPPY_PATH_V2
bash all-up-test.sh VALIDATOR_OUT
# An alchemy api key is necessary for this arbitrary logic test, so we are disabling it
#bash all-up-test.sh ARBITRARY_LOGIC
echo "All tests succeeded!"
