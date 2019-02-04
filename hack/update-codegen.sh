#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPTPATH="$(cd "$(dirname "$0")" && pwd -P)"
GENERATOR_BASE=${SCRIPTPATH}/../vendor/k8s.io/code-generator

"$GENERATOR_BASE/generate-groups.sh" "all" \
    github.com/ctron/iot-simulator-operator/pkg/client \
    github.com/ctron/iot-simulator-operator/pkg/apis \
    "simulator:v1alpha1" \
    --go-header-file "${SCRIPTPATH}/header.txt" \
    "$@"
