resource "aws_security_group" "k3s_sg" {
  name        = "${var.app_name}-${var.environment}-k3s-sg"
  description = "Security group for K3s cluster nodes"
  vpc_id      = var.vpc_id
}

resource "aws_security_group_rule" "allow_http_from_lb" {
  from_port                = 80
  to_port                  = 80
  type                     = "ingress"
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.lb_sg.id
  security_group_id        = aws_security_group.k3s_sg.id
}

resource "aws_security_group_rule" "allow_https_from_lb" {
  from_port                = 443
  to_port                  = 443
  type                     = "ingress"
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.lb_sg.id
  security_group_id        = aws_security_group.k3s_sg.id
}

resource "aws_security_group_rule" "allow_ssh" {
  from_port         = 22
  to_port           = 22
  type              = "ingress"
  protocol          = "tcp"
  cidr_blocks       = [var.admin_ip]
  security_group_id = aws_security_group.k3s_sg.id
}

resource "aws_security_group_rule" "allow_k3s_api_for_admin" {
  from_port         = 6443
  to_port           = 6443
  type              = "ingress"
  protocol          = "tcp"
  cidr_blocks       = [var.admin_ip]
  security_group_id = aws_security_group.k3s_sg.id
}

resource "aws_security_group_rule" "allow_k3s_api_for_nodes" {
  from_port                = 6443
  to_port                  = 6443
  type                     = "ingress"
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.k3s_sg.id
  security_group_id        = aws_security_group.k3s_sg.id
}


resource "aws_security_group_rule" "allow_flannel_xvlan" {
  from_port                = 8472
  to_port                  = 8472
  type                     = "ingress"
  protocol                 = "udp"
  source_security_group_id = aws_security_group.k3s_sg.id
  security_group_id        = aws_security_group.k3s_sg.id
}

resource "aws_security_group_rule" "allow_kubelet_api" {
  from_port                = 10250
  to_port                  = 10250
  type                     = "ingress"
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.k3s_sg.id
  security_group_id        = aws_security_group.k3s_sg.id
}

resource "aws_security_group_rule" "allow_etcd" {
  from_port                = 2379
  to_port                  = 2380
  type                     = "ingress"
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.k3s_sg.id
  security_group_id        = aws_security_group.k3s_sg.id
}

resource "aws_security_group_rule" "allow_egress_for_k3s" {
  from_port         = 0
  to_port           = 0
  type              = "egress"
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.k3s_sg.id
}
