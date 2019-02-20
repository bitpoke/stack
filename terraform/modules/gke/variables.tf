variable "cluster_name" {
  type = "string"
}

variable "region" {
  type = "string"
}

variable "preemptible" {
  default = false
}

variable "system_node_type" {
  type    = "string"
  default = "n1-standard-2"
}

variable "ingress_node_type" {
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

variable "zones" {
  type    = "list"
  default = []
}
