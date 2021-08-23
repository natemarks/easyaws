
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
  test_name = "easyaws_iam_assume_role"
  default_tags = {
    terraform = "true"
    terragrunt = "true"
    test_name = local.test_name
    deleteme = "true"
    project = "easyaws"
  }
}


# ---------------------------------------------------------------------------------------------------------------------
# Create a role with an attached policy that grants access to some test resource
# in this case we're using parameter store
# ---------------------------------------------------------------------------------------------------------------------
data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "role_trust" {

  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }
  }
}

data "aws_iam_policy_document" "role_can_list_s3" {
  statement {
    sid = "1"

    actions = [
      "s3:ListAllMyBuckets",
      "s3:GetBucketLocation",
    ]

    resources = [
      "arn:aws:s3:::*",
    ]
  }
}

resource "aws_iam_policy" "target" {
  name        = "${local.test_name}_target"
  description = "Policy attached to the iam_assume_role target role"
  policy      = "${join("", data.aws_iam_policy_document.role_can_list_s3.*.json)}"

  tags = merge(local.default_tags,
  var.custom_tags)
}

resource "aws_iam_role" "target" {
  name = local.test_name
  tags = merge(local.default_tags,
  var.custom_tags)

  assume_role_policy = "${join("", data.aws_iam_policy_document.role_trust.*.json)}"
  managed_policy_arns = [aws_iam_policy.target.arn]

}




# -----------------------------------------------------------------------------
# Create a group. It's members will be permitted to assume a test role.
# Attaching policies to groups rather than users is a best practice
# -----------------------------------------------------------------------------

data "aws_iam_policy_document" "assume_role_admin" {

  statement {
    actions   = ["sts:AssumeRole"]
    resources = ["${join("", aws_iam_role.target.*.arn)}"]
  }
}

resource "aws_iam_policy" "assume_role_admin" {
  name        = "permit-assume-role"
  description = "Allow assuming admin role"
  policy      = "${join("", data.aws_iam_policy_document.assume_role_admin.*.json)}"
  tags = merge(local.default_tags,
  var.custom_tags)
}

resource "aws_iam_group" "group" {
  name = local.test_name
  path = "/deleteme/groups/"
}

resource "aws_iam_group_policy_attachment" "assume_role_admin" {
  group      = "${join("", aws_iam_group.group.*.name)}"
  policy_arn = "${join("", aws_iam_policy.assume_role_admin.*.arn)}"
}

# -----------------------------------------------------------------------------
# Create an IAM user  and it access key
# Store the access key in a secret
# Attaching policies to groups rather than users is a best practice
# -----------------------------------------------------------------------------

data "aws_iam_user" "user" {
  user_name = "test_${local.default_tags.project}"
}


resource "aws_iam_group_membership" "group" {
  name = local.test_name

  users = [
    "${data.aws_iam_user.user.user_name}"
  ]
  group = aws_iam_group.group.name
}







