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

variable "subnet_config" {
  description = "config for subnets"

  type = map(object({
    cidr_block = string
    az         = string
  }))

  default = {
    "public_1a" = { cidr_block = "10.0.1.0/24", az = "eu-central-1a" }
    "public_1b" = { cidr_block = "10.0.2.0/24", az = "eu-central-1b" }
    "public_1c" = { cidr_block = "10.0.3.0/24", az = "eu-central-1c" }
  }
}