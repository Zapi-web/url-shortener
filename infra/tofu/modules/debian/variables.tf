variable "debian_version" {
  default     = "12"
  type        = string
  description = "Debian version"

  validation {
    condition     = contains(["11", "12", "13"], var.debian_version)
    error_message = "Only 11-13 versions"
  }
}

variable "arm_instance" {
  type        = bool
  description = "Do you use ARM instance"
  default     = false
}