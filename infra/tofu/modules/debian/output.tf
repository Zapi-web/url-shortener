output "debian_version_id" {
  description = "id of debian version"
  value       = data.aws_ami.debian[var.debian_version].id
}