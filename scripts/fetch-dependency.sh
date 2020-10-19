#!/bin/sh
# fetch-dependency.sh
# This script fetch framework dependency repos.
# @version    190625:2
# @author     karminski <code.karminski@outlook.com>
#

# [ci]
echo "-[fetch dependency]-"
echo "start"

# ----------------------------[manual config here]------------------------------

REPO_PATH='../'

# ----- start

export GOPROXY=https://goproxy.io,direct

cd ${REPO_PATH}

go mod vendor

echo "done"


