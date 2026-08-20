output "pg_backup_bucket_arn" {
  value       = aws_s3_bucket.pg_backups.arn
  description = "Backup Bucket ARN for PostgreSQL"
}