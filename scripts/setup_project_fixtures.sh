#!/usr/bin/env bash
set -Eeuo pipefail


# source import utility functions
. scripts/utility.sh

setup_terraform
run_tf_module deployments/test_user