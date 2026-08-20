output "iam_instance_profile" {
  value       = aws_iam_instance_profile.k3s_instance_profile
  description = "IAM instance profile for EBS and Cluster Scale"
}