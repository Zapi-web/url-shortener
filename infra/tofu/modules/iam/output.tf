output "iam_instance_profile" {
  value       = aws_iam_instance_profile.k3s_instance_profile.name
  description = "IAM instance profile for EBS and Cluster Scale"
}