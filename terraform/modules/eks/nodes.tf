data "aws_ami" "worker" {
  filter {
    name   = "name"
    values = ["amazon-eks-node-${aws_eks_cluster.stack.version}-v*"]
  }

  most_recent = true
  owners      = ["602401143452"] # Amazon EKS AMI Account ID
}

module "system_nodes" {
  source               = "modules/eks-node"
  cluster_name         = "${var.cluster_name}"
  pool_name            = "system"
  instance_type        = "${var.system_instance_type}"
  disk_size_gb         = "50"
  iam_instance_profile = "${aws_iam_instance_profile.node.name}"
  security_groups      = ["${aws_security_group.node.id}"]
  image_id             = "${data.aws_ami.worker.id}"
  vpc_zone_identifier  = ["${aws_subnet.node.*.id}"]
  cluster_endpoint     = "${aws_eks_cluster.stack.endpoint}"
  cluster_ca           = "${aws_eks_cluster.stack.certificate_authority.0.data}"
  min_size             = 1
  max_size             = 3
  desired_capacity     = 1
  spot_price           = "${var.system_spot_price}"

  labels = {
    "node-role.stack.presslabs.org/presslabs-sys" = ""
  }

  taints = {
    "CriticalAddonsOnly" = "true:PreferNoSchedule"
  }
}

module "database_nodes" {
  source               = "modules/eks-node"
  cluster_name         = "${var.cluster_name}"
  pool_name            = "database"
  instance_type        = "${var.database_instance_type}"
  disk_size_gb         = "100"
  iam_instance_profile = "${aws_iam_instance_profile.node.name}"
  security_groups      = ["${aws_security_group.node.id}"]
  image_id             = "${data.aws_ami.worker.id}"
  vpc_zone_identifier  = ["${aws_subnet.node.*.id}"]
  cluster_endpoint     = "${aws_eks_cluster.stack.endpoint}"
  cluster_ca           = "${aws_eks_cluster.stack.certificate_authority.0.data}"
  min_size             = 0
  max_size             = 3
  spot_price           = "${var.database_spot_price}"

  labels = {
    "node-role.stack.presslabs.org/database"  = ""
    "node-role.stack.presslabs.org/mysql"     = ""
    "node-role.stack.presslabs.org/memcached" = ""
  }
}

module "wordpress_nodes" {
  source               = "modules/eks-node"
  cluster_name         = "${var.cluster_name}"
  pool_name            = "wordpress"
  instance_type        = "${var.wordpress_instance_type}"
  disk_size_gb         = "100"
  iam_instance_profile = "${aws_iam_instance_profile.node.name}"
  security_groups      = ["${aws_security_group.node.id}"]
  image_id             = "${data.aws_ami.worker.id}"
  vpc_zone_identifier  = ["${aws_subnet.node.*.id}"]
  cluster_endpoint     = "${aws_eks_cluster.stack.endpoint}"
  cluster_ca           = "${aws_eks_cluster.stack.certificate_authority.0.data}"
  min_size             = 0
  max_size             = 5
  spot_price           = "${var.wordpress_spot_only ? var.wordpress_spot_price : ""}"

  labels = {
    "node-role.stack.presslabs.org/wordpress" = ""
  }
}

module "wordpress_spot_nodes" {
  source               = "modules/eks-node"
  cluster_name         = "${var.cluster_name}"
  pool_name            = "wordpress-spot"
  instance_type        = "${var.wordpress_instance_type}"
  disk_size_gb         = "100"
  iam_instance_profile = "${aws_iam_instance_profile.node.name}"
  security_groups      = ["${aws_security_group.node.id}"]
  image_id             = "${data.aws_ami.worker.id}"
  vpc_zone_identifier  = ["${aws_subnet.node.*.id}"]
  cluster_endpoint     = "${aws_eks_cluster.stack.endpoint}"
  cluster_ca           = "${aws_eks_cluster.stack.certificate_authority.0.data}"
  min_size             = 0
  max_size             = 5
  spot_price           = "${var.wordpress_spot_price}"

  labels = {
    "node-role.stack.presslabs.org/wordpress" = ""
  }

  taints = {
    "stack.presslabs.org/spot-node" = "true:NoSchedule"
  }
}
