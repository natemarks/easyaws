#!/usr/bin/env bash
set -Eeuo pipefail

# Version of pipeline scripts to use
# https://github.com/natemarks/pipeline-scripts
declare -r PS_VER=v0.0.16
declare -r TF_VER=1.0.5

#######################################
# Invoke the terraform install script from the github pipeline-scripts project
#######################################
function setup_terraform() {
  curl -sS "https://raw.githubusercontent.com/natemarks/pipeline-scripts/${PS_VER}/scripts/install_terraform.sh" | bash -s -- -d build/terraform -r "${TF_VER}"
  export PATH="$(pwd)/build/terraform/${TF_VER}:$PATH"
}

#######################################
# Assuming the function is called from the integration test script, get the test name from the integration script
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
  terraform init
  terraform apply -auto-approve


}


#######################################
# Given a path to a terraform module, apply it
#######################################
function applyTerraform() {
  local tf_module="${1}"
  setup_terraform
  bash -c "curl -sS https://raw.githubusercontent.com/natemarks/pipeline-scripts/${PS_VER}/scripts/apply_terraform.sh | bash -s --  -t ${tf_module}"

}


#######################################
# Given a path to a terraform module, destroy it
#######################################
function destroyTerraform() {
  local tf_module="${1}"
  setup_terraform
  bash -c "curl -sS https://raw.githubusercontent.com/natemarks/pipeline-scripts/${PS_VER}/scripts/destroy_terraform.sh | bash -s --  -t ${tf_module}"
}

#######################################
# Given a path to the integration test script (ex. ./scripts/test_iam_assume_role.sh)
# find the terraform module for it's test fixtures and clean them up (terraform destroy)
#
# NOTE: This is used from the Makefile loop through integration test scripts
#######################################
function teardownTestFixtures() {
  # example script name: test_iam_assume_role.sh
  local script_name="$(basename ${1})"
  # example test name: iam_assume_role
  local test_name="$(get_test_name ${script_name})"

  local tf_module="deployments/${test_name}"
  export PATH="$(pwd)/build/terraform/${TF_VER}:$PATH"
  bash -c "curl -sS https://raw.githubusercontent.com/natemarks/pipeline-scripts/${PS_VER}/scripts/destroy_terraform.sh | bash -s --  -t ${tf_module}"

}