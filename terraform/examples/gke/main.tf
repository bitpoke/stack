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

module "stack" {
  source       = "../../modules/gke"
  cluster_name = "${var.cluster_name}"
  region       = "${var.region}"
  zones        = "${var.zones}"
  preemptible  = "${var.preemptible}"
}
