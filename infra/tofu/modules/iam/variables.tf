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

variable "pg_s3_bucket_arn" {
  type        = string
  description = "PostgreSQL S3 bucket ARN"
}