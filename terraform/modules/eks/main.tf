data "aws_eks_cluster_auth" "stack" {
  name = "${aws_eks_cluster.stack.name}"
}

resource "aws_eks_cluster" "stack" {
  name     = "${var.cluster_name}"
  role_arn = "${aws_iam_role.master.arn}"

  vpc_config {
    security_group_ids = ["${aws_security_group.master.id}"]
    subnet_ids         = ["${aws_subnet.node.*.id}"]
  }

  depends_on = [
    "aws_iam_role_policy_attachment.stack-AmazonEKSClusterPolicy",
    "aws_iam_role_policy_attachment.stack-AmazonEKSServicePolicy",
  ]
}

provider "kubernetes" {
  host                   = "${aws_eks_cluster.stack.endpoint}"
  cluster_ca_certificate = "${base64decode(aws_eks_cluster.stack.certificate_authority.0.data)}"
  token                  = "${data.aws_eks_cluster_auth.stack.token}"
  load_config_file       = false
}

resource "kubernetes_config_map" "aws_auth" {
  metadata {
    name      = "aws-auth"
    namespace = "kube-system"
  }

  data {
    mapRoles = <<YAML
- rolearn: ${aws_iam_role.node.arn}
  username: system:node:{{EC2PrivateDNSName}}
  groups:
    - system:bootstrappers
    - system:nodes
YAML
  }
}
