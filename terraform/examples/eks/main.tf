provider "aws" {
  version = "~> 2.0"
  region  = "${var.region}"
}

variable "cluster_name" {
  type = "string"
}

variable "region" {
  type    = "string"
  default = "eu-central-1"
}

variable "zones" {
  type    = "list"
  default = ["eu-central-1a", "eu-central-1b"]
}

variable "system_instance_type" {
  type    = "string"
  default = "m5.large"
}

variable "database_instance_type" {
  type    = "string"
  default = "m5.xlarge"
}

variable "wordpress_instance_type" {
  type    = "string"
  default = "m5.xlarge"
}

variable "wordpress_spot_price" {
  type    = "string"
  default = "0.14"
}

variable "system_spot_price" {
  type    = "string"
  default = ""
}

variable "database_spot_price" {
  type    = "string"
  default = ""
}

variable "wordpress_spot_only" {
  default = false
}

module "stack" {
  source       = "../../modules/eks"
  cluster_name = "${var.cluster_name}"
  zones        = "${var.zones}"

  system_instance_type = "${var.system_instance_type}"
  system_spot_price    = "${var.system_spot_price}"

  database_instance_type = "${var.database_instance_type}"
  database_spot_price    = "${var.database_spot_price}"

  wordpress_instance_type = "${var.wordpress_instance_type}"
  wordpress_spot_price    = "${var.wordpress_spot_price}"
  wordpress_spot_only     = "${var.wordpress_spot_only}"
}
