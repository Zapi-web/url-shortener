locals {
  arch_name = var.arm_instance ? "arm64" : "amd64"
  arch_ec2  = var.arm_instance ? "arm64" : "x86_64"
}

data "aws_ami" "debian" {
  for_each = toset(["11", "12", "13"])

  most_recent = true
  owners      = ["136693071363"]

  filter {
    name   = "name"
    values = ["debian-${each.value}-${local.arch_name}-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = [local.arch_ec2]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
}