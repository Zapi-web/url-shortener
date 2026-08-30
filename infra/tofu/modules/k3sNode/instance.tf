resource "aws_instance" "k3s_node_instance" {
  for_each = var.subnet_ids

  ami                    = var.debian_version_data_id
  instance_type          = var.instance_type
  iam_instance_profile   = var.k3s_iam_instance_profile_name
  subnet_id              = each.value
  vpc_security_group_ids = [var.k3s_nodes_sg_id]
  key_name               = var.key_name

  root_block_device {
    volume_type = "gp3"
    volume_size = 20
  }

  tags = {
    Name        = "${var.app_name}-${var.environment}-node-instance-${each.key}"
    Environment = var.environment
    Role        = "server"
    Cluster     = "${var.app_name}-${var.environment}-k3s"
  }
}