variable "app_name" {
  type        = string
  default     = "default-app"
  description = "App name"
}

variable "environment" {
  type        = string
  default     = "test"
  description = "App environment"

  validation {
    condition     = contains(["dev", "test", "prod"], var.environment)
    error_message = "Environment must be dev,test or prod"
  }
}

variable "vpc_id" {
  type        = string
  description = "ID of VPC"
}

variable "subnet_ids" {
  type        = map(string)
  description = "IDs of public subnets"
}

variable "k3s_nodes_ids" {
  type        = map(string)
  description = "IDs of k3s nodes"
}

variable "lb_sg_id" {
  type        = string
  description = "Load Balancer security group id"
}

variable "lb_listening_port" {
  type        = number
  description = "What port should LB listen"
}