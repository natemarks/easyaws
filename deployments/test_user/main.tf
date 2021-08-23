
# -----------------------------------------------------------------------------
# This module manages the iam_assume_role integration test fixture
# The test requires:
# - an IAM user in a group that has rights to assume a target role
# - an assumable role that has policy access to some test resource
# - a resource that can be accessed only by assuming the target role
# -----------------------------------------------------------------------------


provider "aws" {
  region = "us-east-1"
}

terraform {
  required_version = ">= 1.0.5"

  required_providers {
    aws = ">= 3.0.0"
  }
}

locals {
  default_tags = {
    terraform = "true"
    terragrunt = "true"
    deleteme = "true"
    project = "easyaws"
  }
}


# -----------------------------------------------------------------------------
# Create an IAM user  and it access key
# Store the access key in a secret
# Attaching policies to groups rather than users is a best practice
# -----------------------------------------------------------------------------

resource "aws_iam_user" "user" {
  name = "test_${local.default_tags.project}"
  path = "/deleteme/users/"
  tags = merge(local.default_tags,
  var.custom_tags)
}

resource "aws_iam_access_key" "user" {
  user    = aws_iam_user.user.name
}

resource aws_secretsmanager_secret "user" {
  name = "${aws_iam_user.user.name}_credentials"
  recovery_window_in_days = 0
  tags = merge(local.default_tags,
  var.custom_tags)
}

resource "aws_secretsmanager_secret_version" "user" {
  secret_id     = "${aws_secretsmanager_secret.user.id}"
  secret_string = jsonencode({"AWSAccessKeyID" = aws_iam_access_key.user.id, "AWSSecretAccessKey" = aws_iam_access_key.user.secret})
}






