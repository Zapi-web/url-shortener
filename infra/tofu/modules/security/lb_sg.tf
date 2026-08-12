resource "aws_security_group" "lb_sg" {
  name        = "${var.app_name}-${var.environment}-lb-sg"
  description = "Security group for Load Balancer"
  vpc_id      = var.vpc_id
}

resource "aws_security_group_rule" "allow_http" {
  from_port         = 80
  to_port           = 80
  type              = "ingress"
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.lb_sg.id
}

resource "aws_security_group_rule" "allow_https" {
  from_port         = 443
  to_port           = 443
  type              = "ingress"
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.lb_sg.id
}

resource "aws_security_group_rule" "allow_egress_for_lb" {
  from_port         = 0
  to_port           = 0
  type              = "egress"
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.lb_sg.id
}
