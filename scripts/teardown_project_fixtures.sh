#!/usr/bin/env bash
set -Eeuo pipefail


# source import utility functions
. scripts/utility.sh

setup_terraform
destroyTerraform deployments/test_user