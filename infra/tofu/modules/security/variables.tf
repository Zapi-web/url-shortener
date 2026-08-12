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
  description = "Admin ip for SSH"
}

variable "vpc_id" {
  description = "ID of VPC"
  type        = string
}