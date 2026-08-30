module "vpc" {
  source = "./modules/vpc"

  app_name    = var.app_name
  environment = var.environment
}

module "debian" {
  source = "./modules/debian"

  arm_instance   = var.arm_instance
  debian_version = var.debian_version
}

module "s3" {
  source = "./modules/s3"

  app_name    = var.app_name
  environment = var.environment
}

module "iam" {
  source = "./modules/iam"

  app_name         = var.app_name
  environment      = var.environment
  pg_s3_bucket_arn = module.s3.pg_backup_bucket_arn
}

module "security_groups" {
  source = "./modules/security"

  app_name    = var.app_name
  environment = var.environment
  vpc_id      = module.vpc.vpc_id
  admin_ip    = var.admin_ip
}

module "k3s_nodes" {
  source = "./modules/k3sNode"

  app_name                      = var.app_name
  environment                   = var.environment
  instance_type                 = var.instances_type
  debian_version_data_id        = module.debian.debian_version_id
  subnet_ids                    = module.vpc.public_subnet_ids
  k3s_nodes_sg_id               = module.security_groups.k3s_nodes_sg_id
  key_name                      = var.key_name
  k3s_iam_instance_profile_name = module.iam.iam_instance_profile
}

module "load_balancer" {
  source = "./modules/loadBalancer"

  app_name          = var.app_name
  environment       = var.environment
  vpc_id            = module.vpc.vpc_id
  subnet_ids        = module.vpc.public_subnet_ids
  k3s_nodes_ids     = module.k3s_nodes.k3s_nodes_ids
  lb_sg_id          = module.security_groups.lb_sg_id
  lb_listening_port = var.lb_listening_port
}