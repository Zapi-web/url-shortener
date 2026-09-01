resource "aws_s3_bucket" "pg_backups" {
  bucket = "${var.app_name}-${var.environment}-pg-backups"

  lifecycle {
    prevent_destroy = false
  }
}

resource "aws_s3_bucket_public_access_block" "pg_backups" {
  bucket = aws_s3_bucket.pg_backups.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_lifecycle_configuration" "pg_backups" {
  bucket = aws_s3_bucket.pg_backups.id

  rule {
    id     = "cleanup-old-backups"
    status = "Enabled"

    expiration {
      days = 14
    }
  }
}