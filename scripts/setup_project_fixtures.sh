#!/usr/bin/env bash
set -Eeuo pipefail


# source import utility functions
. scripts/utility.sh

setup_terraform
applyTerraform deployments/test_user