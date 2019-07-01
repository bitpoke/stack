variable "cluster_name" {
  type        = "string"
  description = "GKE cluster name"
}

variable "region" {
  type        = "string"
  description = "The GCP region to start your cluster in"
}

variable "project" {
  type    = "string"
  default = ""
}

variable "preemptible" {
  default     = false
  description = <<EOF
    Whenever to create a preemptible-only cluster. By default, the pools
    running WordPress sites are split into preemptible/non-preemptible nodes
    and the preemptible nodes are
    tainted with cloud.google.com/gke-preemptible=true:NoSchedule
  EOF
}

variable "system_node_taint_effect" {
  type = "string"
  default = "NO_SCHEDULE"
  description = <<EOF
    The system node pool is tainted with CriticalAddonsOnly.
    This is the default effect for that taint.
  EOF
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

variable "zones" {
  type    = "list"
  default = []
}
