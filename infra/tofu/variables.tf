variable "app_name" {
  type        = string
  description = "App name"
  default     = "default-app"
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

variable "debian_version" {
  default     = "12"
  type        = string
  description = "Debian version"

  validation {
    condition     = contains(["11", "12", "13"], var.debian_version)
    error_message = "Only 11-13 versions"
  }
}

variable "admin_ip" {
  type        = string
  description = "Admin ip for SSH" # CAREFUL! Must be CIDR notaion "YOUR_IP/32"

  validation {
    condition     = can(cidrnetmask(var.admin_ip))
    error_message = "admin_ip must be a valid CIDR string, e.g., '203.0.113.5/32'."
  }
}

variable "arm_instance" {
  type        = bool
  description = "Do you use ARM instance"
  default     = false
}

variable "instances_type" {
  type        = string
  description = "Instance type"
  default     = "t3.medium"
}

variable "key_name" {
  type        = string
  description = "SSH key"
}

variable "lb_listening_port" {
  type        = number
  description = "Load Balancer listening port"
  default     = 80
}

variable "aws_region" {
  type        = string
  default     = "eu-central-1"
  description = "AWS Region"
}