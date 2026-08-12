data "aws_ami" "debian" {
  for_each = toset(["10", "11", "12", "13"])

  most_recent = true
  owners      = ["136693071363"]

  filter {
    name   = "name"
    values = ["debian-${each.value}-amd64-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
}