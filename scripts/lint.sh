#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname "$0")"

pushd "$SCRIPT_DIR/.." > /dev/null || exit 1

    APP_NAME="$(basename "$(pwd)")"
    PACKAGES=$(go list ./... | grep -v mocks | sed -e "s/bitbucket\\.org\\/icelolly\\/$APP_NAME/\\./")

    golint -set_exit_status ${PACKAGES}
    go vet ${PACKAGES}

    # TODO: Replace this.
    #gometalinter --fast --errors --deadline=60s --disable=gotype ${PACKAGES}

popd > /dev/null || exit 1
