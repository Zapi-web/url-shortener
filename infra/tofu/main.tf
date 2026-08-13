module "vpc" {
  source = "./modules/vpc"

  app_name    = var.app_name
  environment = var.environment
}

module "debian" {
  source = "./modules/debian"

  debian_version = var.debian_version
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

  app_name               = var.app_name
  environment            = var.environment
  instance_type          = var.instances_type
  debian_version_data_id = module.debian.debian_version_id
  subnet_ids             = module.vpc.public_subnet_ids
  k3s_nodes_sg_id        = module.security_groups.k3s_nodes_sg_id
  key_name               = var.key_name
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