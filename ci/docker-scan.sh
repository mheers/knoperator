#!/usr/bin/env bash

# BEWARE!
# This script is managed in pipeline scripts repo and copied by GitLab CI in all the other repos.

set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" ; pwd -P)"
. "${SCRIPT_DIR}/helper.inc.sh"

function usage() {
  echo "usage: $0 --registry <REGISTRY>"
  echo "        --registry, -r: set the docker registry"
  echo "        --commit-sha, -c: set the commit-sha"
}

while [[ $# -gt 0 ]]; do
  key="$1"

  case $key in
  --registry|-r)
    export DOCKER_REGISTRY="$2"
    shift
    shift
    ;;
  --commit-sha|-c)
    export COMMIT_SHA="$2"
    shift
    shift
    ;;
  --help|help|-h)
    usage
    exit 0
    shift
    ;;
  *)
    shift
    shift
    ;;
  esac
done

if [[ -z "$DOCKER_REGISTRY" ]]; then
  usage
  echo
  echo "Must pass registry as parameter. Exiting..."
  exit 1
fi

IMAGE_ID="${DOCKER_REGISTRY}/${SERVICE_NAME}:${COMMIT_SHA}"

print-banner "Scanning ${IMAGE_ID}"

if [[ "$OSTYPE" == "darwin"* ]]; then
  OS_RELEASE="macOS"
else
  OS_RELEASE="Linux"
fi
TRIVY_VERSION=$(curl --silent "https://api.github.com/repos/aquasecurity/trivy/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
echo "Using trivy version ${TRIVY_VERSION}"
wget -q "https://github.com/aquasecurity/trivy/releases/download/v${TRIVY_VERSION}/trivy_${TRIVY_VERSION}_${OS_RELEASE}-64bit.tar.gz" -O "/tmp/trivy.tar.gz"
pushd /tmp > /dev/null
tar zxf "trivy.tar.gz"

./trivy image --exit-code 0 --no-progress --ignore-unfixed --severity HIGH "${IMAGE_ID}"
./trivy image --exit-code 1 --severity CRITICAL --no-progress --ignore-unfixed "${IMAGE_ID}"