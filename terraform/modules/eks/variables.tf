variable "cluster_name" {
  type = "string"
}

variable "zones" {
  type = "list"
}

variable "system_instance_type" {
  default = "m5.large"
}

variable "system_spot_price" {
  type    = "string"
  default = ""
}

variable "database_instance_type" {
  default = "m5.xlarge"
}

variable "database_spot_price" {
  type    = "string"
  default = ""
}

variable "wordpress_spot_price" {
  default = "0.14"
}

variable "wordpress_instance_type" {
  default = "m5.xlarge"
}

variable "wordpress_spot_only" {
  default = false
}
