#!/usr/bin/env bash
set -ex

VERSION=$(git describe --always --tags --long)
PLATFORM=""

if [[ ${RUNNER_OS} == 'Linux' ]]; then
  PLATFORM="linux"
elif [[ ${RUNNER_OS} == 'macOS' ]]; then
  PLATFORM="darwin"
else
  PLATFORM="windows"
  exit 1
fi



env GO111MODULE=on ../scripts/build-abigen.sh
