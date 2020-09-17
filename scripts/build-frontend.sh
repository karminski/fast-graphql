#!/bin/bash
# build-frontend.sh
# build repo.
# @version    170322:1
# @author     karminski <code.karminski@outlook.com>
#

# ----------------------------[manual config here]------------------------------

REPO_PATH='/data/repo/jinkanhq/fast-graphql'



# ------------------------------------------------------------------------------

cd "${REPO_PATH}/src/cmd/fast-graphql-frontend/"
go build -mod=vendor -o ${REPO_PATH}/bin/fast-graphql-frontend



exit 0;
