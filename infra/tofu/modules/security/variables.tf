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

variable "admin_ip" {
  type        = string
  description = "Admin ip for SSH" # CAREFUL! Must be CIDR notaion "YOUR_IP/32"

  validation {
    condition     = can(cidrnetmask(var.admin_ip))
    error_message = "admin_ip must be a valid CIDR string, e.g., '203.0.113.5/32'."
  }
}

variable "vpc_id" {
  description = "ID of VPC"
  type        = string
}