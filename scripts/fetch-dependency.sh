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

REPO_PATH='/data/repo/jinkanhq/fast-graphql'

# ----- start

cd ${REPO_PATH}

go mod vendor

echo "done"


