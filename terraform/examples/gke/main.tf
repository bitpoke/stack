provider "google-beta" {
  version = "~> 2.0"
  project = "${var.project}"
  region  = "${var.region}"
  zone    = "${var.zone}"
}

variable "project" {
  type = "string"
}

variable "region" {
  type    = "string"
  default = "europe-west3"
}

variable "zone" {
  type    = "string"
  default = "europe-west3-a"
}

variable "zones" {
  type    = "list"
  default = ["europe-west3-a", "europe-west3-b"]
}

variable "preemptible" {
  default = false
}

variable "cluster_name" {
  type = "string"
}

variable "system_node_type" {
  type    = "string"
  default = "n1-standard-2"
}

variable "database_node_type" {
  type    = "string"
  default = "n1-standard-2"
}

variable "wordpress_node_type" {
  type    = "string"
  default = "n1-standard-4"
}

module "stack" {
  source              = "../../modules/gke"
  cluster_name        = "${var.cluster_name}"
  project             = "${var.project}"
  region              = "${var.region}"
  zones               = "${var.zones}"
  preemptible         = "${var.preemptible}"
  system_node_type    = "${var.system_node_type}"
  database_node_type  = "${var.database_node_type}"
  wordpress_node_type = "${var.wordpress_node_type}"
}
