output "k3s_nodes_sg_id" {
  value       = aws_security_group.k3s_sg.id
  description = "K3s nodes security group id"
}

output "lb_sg_id" {
  value       = aws_security_group.lb_sg.id
  description = "Load Balancer security group id"
}