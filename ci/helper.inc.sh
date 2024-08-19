# BEWARE!
# This script is managed in pipeline scripts repo and copied by GitLab CI in all the other repos.

. "${SCRIPT_DIR}/config.global.inc.sh"
. "${SCRIPT_DIR}/config.service.inc.sh"


function print-banner() {
  echo "#########################################################################################################"
  echo "# ${*} #"
  echo "#########################################################################################################"
}

function log() {
  echo "$(date) ${*}"
}

