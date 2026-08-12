output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.vpc.id
}

output "public_subnet_ids" {
  description = "Map of IDs of public subnets"
  value       = { for key, subnet in aws_subnet.public_subnet : key => subnet.id }
}