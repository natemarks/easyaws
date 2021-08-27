#!/usr/bin/env bash
set -Eeuo pipefail

# Version of pipeline scripts to use
# https://github.com/natemarks/pipeline-scripts
declare -r PS_VER=v0.0.13

#######################################
# Invoke the terraform install script from the github pipeline-scripts project
#######################################
function setup_terraform() {
  curl -s "https://raw.githubusercontent.com/natemarks/pipeline-scripts/${PS_VER}/scripts/install_terraform.sh" | bash -s -- -d build/terraform -r 1.0.5
  export PATH="$(pwd)/build/terraform/1.0.5:$PATH"
}

#######################################
# Assuming the funciton is called from the integration test script, get the test name from the integration script
# name.  This is required to invoke terraform against the module based on the test name.
# Example, the script test_iam_assume_role.sh should return gthe test name 'iam_assume_role', which allows the script
# to invoke terraform in deployments/iam_assume_role
#######################################
function get_test_name() {
    local script_name="${1}"
    local no_suffix="${script_name%.sh}"
    echo "${no_suffix#test_}"
}


#######################################
# Change direc
# name.  This is required to invoke terraform against the module based on the test name.
# Example, the script test_iam_assume_role.sh should return gthe test name 'iam_assume_role', which allows the script
# to invoke terraform in deployments/iam_assume_role
#######################################
function run_tf_module() {
  local module_dir="${1}"
  local initial_dir="$(pwd)"
  cd "${module_dir}"
  terraform init && terraform plan
  terraform apply -auto-approve


}