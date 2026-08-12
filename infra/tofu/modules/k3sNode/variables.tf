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

variable "instance_type" {
  type        = string
  description = "Instance type for app-server"
  default     = "t3.micro"
}

variable "debian_version_data_id" {
  type        = string
  description = "Data id of debian version server"
}

variable "subnet_ids" {
  type        = map(string)
  description = "IDs of public subnets"
}

variable "k3s_nodes_sg_id" {
  type        = string
  description = "K3s nodes security group ID"
}

variable "key_name" {
  type        = string
  description = "SSH key"
}