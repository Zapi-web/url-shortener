variable "debian_version" {
  default     = "12"
  type        = string
  description = "Debian version"

  validation {
    condition     = contains(["10", "11", "12", "13"], var.debian_version)
    error_message = "Only 10-13 versions"
  }
}