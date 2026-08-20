resource "aws_iam_role" "k3s_node_iam_role" {
  name = "k3s-${var.app_name}-${var.environment}-cluster-node-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })
}

resource "aws_iam_policy" "k3s_ebs_csi" {
  name        = "k3s-${var.app_name}-${var.environment}-ebs-csi-policy"
  description = "Policy for EBS-drivers in K3s cluster"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ec2:CreateVolume",
          "ec2:DeleteVolume",
          "ec2:AttachVolume",
          "ec2:DetachVolume",
          "ec2:ModifyVolume",
          "ec2:DescribeVolumes",
          "ec2:DescribeVolumeStatus",
          "ec2:DescribeAvailabilityZones",
          "ec2:DescribeInstances",
          "ec2:DescribeSnapshots",
          "ec2:CreateSnapshot",
          "ec2:DeleteSnapshot",
          "ec2:CreateTags"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_policy" "k3s_s3_backup" {
  name        = "${var.app_name}-${var.environment}-s3-policy"
  description = "Policy for S3 PG backups"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ],
        Resource = [
          var.pg_s3_bucket_arn,
          "${var.pg_s3_bucket_arn}/*"
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ebs_csi_attach" {
  role       = aws_iam_role.k3s_node_iam_role.name
  policy_arn = aws_iam_policy.k3s_ebs_csi.arn
}

resource "aws_iam_role_policy_attachment" "s3_attach" {
  role       = aws_iam_role.k3s_node_iam_role.name
  policy_arn = aws_iam_policy.k3s_s3_backup.arn
}

resource "aws_iam_instance_profile" "k3s_instance_profile" {
  name = "k3s-${var.app_name}-${var.environment}-instance-profile"
  role = aws_iam_role.k3s_node_iam_role.name
}