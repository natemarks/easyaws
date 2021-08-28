#!/usr/bin/env bash
# https://gist.github.com/natemarks/aebb7e84010d4bc37270d554106cb38b
set -Eeuo pipefail
trap cleanup SIGINT SIGTERM ERR EXIT

# shellcheck disable=SC2034
script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)

usage() {
  cat <<EOF
Usage: $(basename "${BASH_SOURCE[0]}") [-h] [-v]

Script description here.

Available options:

-h, --help      Print this help and exit
-v, --verbose   Print script debug info
EOF
  exit
}

cleanup() {
  trap - SIGINT SIGTERM ERR EXIT
  setup_colors
  msg "${RED}CLEAN UP: Running terraform destroy${NOFORMAT}"

}

setup_colors() {
  if [[ -t 2 ]] && [[ -z "${NO_COLOR-}" ]] && [[ "${TERM-}" != "dumb" ]]; then
    NOFORMAT='\033[0m' RED='\033[0;31m' GREEN='\033[0;32m' ORANGE='\033[0;33m' BLUE='\033[0;34m' PURPLE='\033[0;35m' CYAN='\033[0;36m' YELLOW='\033[1;33m'
  else
    # shellcheck disable=SC2034
    NOFORMAT='' RED='' GREEN='' ORANGE='' BLUE='' PURPLE='' CYAN='' YELLOW=''
  fi
}

msg() {
  echo >&2 -e "${1-}"
}

die() {
  local msg=$1
  local code=${2-1} # default exit status 1
  msg "$msg"
  exit "$code"
}

parse_params() {


  while :; do
    case "${1-}" in
    -h | --help) usage ;;
    -v | --verbose) set -xv ;;
    --no-color) NO_COLOR=1 ;;
    -?*) die "Unknown option: $1" ;;
    *) break ;;
    esac
    shift
  done

  args=("$@")

  return 0
}

parse_params "$@"
setup_colors
# source import utility functions
. scripts/utility.sh
# install terraform and add it to the path
setup_terraform
# identify the terraform module for the resources this test requires
declare -r tf_module="deployments/$(get_test_name "$(basename ${BASH_SOURCE[0]})")"
# NOTE PS_VER is set by sourcing utility.sh. This lets me manage that project dependency in a single place
source /dev/stdin <<<"$( curl -sS https://raw.githubusercontent.com/natemarks/pipeline-scripts/${PS_VER}/scripts/utility.sh )"
# TF apply for test resources
bash -c "curl https://raw.githubusercontent.com/natemarks/pipeline-scripts/${PS_VER}/scripts/apply_terraform.sh | bash -s --  -t ${tf_module}"
# export the creds for the project test account
credsFromSecretManager test_easyaws_credentials
aws sts get-caller-identity
# run the tests
sleep 10
go test github.com/natemarks/easyaws/internal --tags=integration
