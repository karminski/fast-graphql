#!/bin/bash
# build.sh
# build repo.
# @version    20200918:1
# @author     karminski <code.karminski@outlook.com>
#

# ----------------------------[manual config here]------------------------------

REPO_PATH='/data/repo/jinkanhq/fast-graphql'



# ------------------------------------------------------------------------------

cd "${REPO_PATH}/src/cmd/fast-graphql/"
go build -mod=vendor -o ${REPO_PATH}/bin/fast-graphql



exit 0;
