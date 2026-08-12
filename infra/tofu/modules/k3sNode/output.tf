output "k3s_nodes_ids" {
  description = "K3s_nodes_ids"
  value       = { for key, instance in aws_instance.k3s_node_instance : key => instance.id }
}