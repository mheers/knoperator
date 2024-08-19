#!/usr/bin/env bash

# BEWARE!
# This script is managed in pipeline scripts repo and copied by GitLab CI in all the other repos.

set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" ; pwd -P)"
. "${SCRIPT_DIR}/helper.inc.sh"
. "${SCRIPT_DIR}/../.VERSION"

function usage() {
  echo "usage: $0 --registry <REGISTRY> [--dry-run]"
  echo "        --registry, -r: set the docker registry"
  echo "        --commit-sha, -c: set the commit-sha"
  echo "        --platform, -p: set the target platform architecture (e.g.  --platform linux/arm64,linux/amd64)"
  echo "        --dry-run:      omit 'docker push'"
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
  --platform|-p)
    export PLATFORM="$2"
    shift
    shift
    ;;
  --dry-run)
    export DRY_RUN="true"
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

if [[ -z "$PLATFORM" ]]; then
  usage
  echo
  echo "Must pass platform as parameter. Exiting..."
  exit 1
fi

IMAGE_ID="${DOCKER_REGISTRY}/${SERVICE_NAME}"

print-banner "Building ${IMAGE_ID} with version ${VERSION} and commit-sha ${COMMIT_SHA} for platforms ${PLATFORM}"

pushd "${SCRIPT_DIR}/.." > /dev/null

PUSH=""
if [[ -z "${DRY_RUN}" ]]; then
  PUSH="--push"
fi
docker buildx build --platform ${PLATFORM} -t "${IMAGE_ID}:${VERSION}" -t "${IMAGE_ID}:${COMMIT_SHA}" ${PUSH} .

exec "${@}"
