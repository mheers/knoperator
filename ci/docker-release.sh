#!/usr/bin/env bash

# BEWARE!
# This script is managed in pipeline scripts repo and copied by GitLab CI in all the other repos.

set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" ; pwd -P)"
. "${SCRIPT_DIR}/helper.inc.sh"
. "${SCRIPT_DIR}/../.VERSION"

function usage() {
  echo "usage: $0 --source-registry <SOURCE> --target-registry <TARGET> [--dry-run]"
  echo "        --source-registry, -s: ze source!"
  echo "        --target-registry, -t: ze target?"
  echo "        --dry-run:          omit 'docker push'"
}

while [[ $# -gt 0 ]]; do
  key="$1"

  case $key in
  --source-registry)
    export SOURCE="$2"
    shift
    shift
    ;;
  --target-registry)
    export TARGET="$2"
    shift
    shift
    ;;
  --dry-run)
    export DRY_RUN="echo"
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

SOURCE_IMAGE="${SOURCE}/${SERVICE_NAME}:${VERSION}"
TARGET_IMAGE="${TARGET}/${SERVICE_NAME}:${VERSION}"
TARGET_IMAGE_LATEST="${TARGET}/${SERVICE_NAME}:latest"

docker pull "${SOURCE_IMAGE}"

docker tag "${SOURCE_IMAGE}" "${TARGET_IMAGE}"
docker tag "${SOURCE_IMAGE}" "${TARGET_IMAGE_LATEST}"

docker push "${TARGET_IMAGE}"
${DRY_RUN} docker push "${TARGET_IMAGE_LATEST}"
