resource "aws_instance" "k3s_node_instance" {
  for_each = var.subnet_ids

  ami                    = var.debian_version_data_id
  instance_type          = var.instance_type
  subnet_id              = each.value
  vpc_security_group_ids = [var.k3s_nodes_sg_id]
  key_name               = var.key_name

  tags = {
    Name        = "${var.app_name}-${var.environment}-node-instance-${each.key}"
    Environment = var.environment
    Role        = "node"
  }
}