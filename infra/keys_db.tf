
resource "aws_vpc" "main" {
  cidr_block = "${var.vpc_cidr_block}"
  enable_dns_support = "true"
  enable_dns_hostnames = "true"
}

resource "aws_subnet" "main_az1" {
  cidr_block = "${lookup(var.main_subnet_cidr_blocks, "0")}"
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "${lookup(var.aws_availability_zones, "0")}"
}

resource "aws_subnet" "main_az2" {
  cidr_block = "${lookup(var.main_subnet_cidr_blocks, "1")}"
  vpc_id = "${aws_vpc.main.id}"
  availability_zone = "${lookup(var.aws_availability_zones, "1")}"
}

resource "aws_internet_gateway" "gw" {
  vpc_id = "${aws_vpc.main.id}"
}

resource "aws_route_table" "vpc" {
  vpc_id = "${aws_vpc.main.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.gw.id}"
  }
}

resource "aws_route_table_association" "main_az1" {
  route_table_id = "${aws_route_table.vpc.id}"
  subnet_id = "${aws_subnet.main_az1.id}"
}

resource "aws_route_table_association" "main_az2" {
  route_table_id = "${aws_route_table.vpc.id}"
  subnet_id = "${aws_subnet.main_az2.id}"
}

resource "aws_db_subnet_group" "default" {
  name = "main"
  subnet_ids = ["${aws_subnet.main_az1.id}", "${aws_subnet.main_az2.id}"]
}

resource "aws_security_group" "keys_db" {
  name = "keys_db"
  vpc_id = "${aws_vpc.main.id}"

  ingress {
    from_port = "${var.db_port}"
    protocol = "tcp"
    to_port = "${var.db_port}"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_instance" "keys" {
  allocated_storage = 10
  storage_type = "gp2"
  engine = "postgres"
  instance_class = "db.t2.micro"
  apply_immediately = "${var.db_unscheduled_restarts_allowed}"
  name = "${var.db_prefix}keys"
  username = "${var.db_username}"
  password = "${var.db_password}"
  db_subnet_group_name = "${aws_db_subnet_group.default.name}"
  multi_az = "false"
  publicly_accessible = "${var.db_public_access}"
  vpc_security_group_ids = ["${aws_security_group.keys_db.id}"]
  port = "${var.db_port}"
  skip_final_snapshot = "true"
  depends_on = ["aws_internet_gateway.gw"]
}

output "keys_db_address" {
  value = "${aws_db_instance.keys.address}"
}

output "keys_db_port" {
  value = "${aws_db_instance.keys.port}"
}

output "keys_db_master_username" {
  value = "${aws_db_instance.keys.username}"
}

output "keys_db_database_name" {
  value = "${aws_db_instance.keys.name}"
}
