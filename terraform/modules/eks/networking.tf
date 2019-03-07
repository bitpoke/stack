resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = "${
    map(
     "Name", "${var.cluster_name}",
     "kubernetes.io/cluster/${var.cluster_name}", "shared",
    )
  }"
}

resource "aws_subnet" "node" {
  count = "${length(var.zones)}"

  availability_zone = "${var.zones[count.index]}"
  cidr_block        = "10.0.${count.index}.0/24"
  vpc_id            = "${aws_vpc.main.id}"

  tags = "${
    map(
     "Name", "${var.cluster_name}-node",
     "kubernetes.io/cluster/${var.cluster_name}", "shared",
    )
  }"
}

resource "aws_internet_gateway" "main" {
  vpc_id = "${aws_vpc.main.id}"

  tags = {
    Name = "${var.cluster_name}"
  }
}

resource "aws_route_table" "main" {
  vpc_id = "${aws_vpc.main.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.main.id}"
  }
}

resource "aws_route_table_association" "main" {
  count = "${length(var.zones)}"

  subnet_id      = "${aws_subnet.node.*.id[count.index]}"
  route_table_id = "${aws_route_table.main.id}"
}
