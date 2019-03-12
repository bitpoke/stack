variable "cluster_name" {
  type = "string"
}

variable "pool_name" {
  type = "string"
}

variable "cluster_endpoint" {
  type = "string"
}

variable "cluster_ca" {
  type = "string"
}

variable "instance_type" {
  type    = "string"
  default = "m5.large"
}

variable "image_id" {
  type = "string"
}

variable "iam_instance_profile" {
  type = "string"
}

variable "security_groups" {
  type = "list"
}

variable "vpc_zone_identifier" {
  type = "list"
}

variable "min_size" {
  default = 0
}

variable "max_size" {
  default = 0
}

variable "desired_capacity" {
  default = 0
}

variable "labels" {
  default = {}
}

variable "taints" {
  default = {}
}

variable "spot_price" {
  default = ""
}

variable "disk_size_gb" {
  default = "100"
}

# EKS currently documents this required userdata for EKS worker nodes to
# properly configure Kubernetes applications on the EC2 instance.
# We utilize a Terraform local here to simplify Base64 encoding this
# information into the AutoScaling Launch Configuration.
# More information: https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
locals {
  node_userdata = <<USERDATA
#!/bin/bash
set -o xtrace
/etc/eks/bootstrap.sh \
  --apiserver-endpoint '${var.cluster_endpoint}' \
  --b64-cluster-ca '${var.cluster_ca}' \
  --kubelet-extra-args "--node-labels=${join(",",formatlist("%s=%s", keys(var.labels), values(var.labels)))} --register-with-taints=${join(",",formatlist("%s=%s", keys(var.taints), values(var.taints)))}" \
  '${var.cluster_name}'
USERDATA
}

resource "aws_launch_configuration" "node" {
  associate_public_ip_address = true
  iam_instance_profile        = "${var.iam_instance_profile}"
  instance_type               = "${var.instance_type}"
  name_prefix                 = "${var.cluster_name}-${var.pool_name}-"
  user_data_base64            = "${base64encode(local.node_userdata)}"
  security_groups             = ["${var.security_groups}"]
  image_id                    = "${var.image_id}"
  spot_price                  = "${var.spot_price}"

  root_block_device = {
    volume_size = "${var.disk_size_gb}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "demo" {
  count                = "${length(var.vpc_zone_identifier)}"
  launch_configuration = "${aws_launch_configuration.node.id}"
  desired_capacity     = "${var.desired_capacity}"
  max_size             = "${var.max_size}"
  min_size             = "${var.min_size}"
  name_prefix          = "${var.cluster_name}-${var.pool_name}-"
  vpc_zone_identifier  = ["${var.vpc_zone_identifier[count.index]}"]

  tag {
    key                 = "Name"
    value               = "${var.cluster_name}-${var.pool_name}"
    propagate_at_launch = true
  }

  tag {
    key                 = "k8s.io/cluster-autoscaler/enabled"
    value               = "true"
    propagate_at_launch = true
  }

  tag {
    key                 = "kubernetes.io/cluster/${var.cluster_name}"
    value               = "owned"
    propagate_at_launch = true
  }
}
