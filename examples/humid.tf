resource "humid" "server" {
    keepers = {
        # Generate a new id each time we switch to a new AMI id
        ami_id = "${var.ami_id}"
    }

    list = "animals"
    adjectives = 2
    separator = "-"
    capitalize = false
}

resource "aws_instance" "server" {
    tags = {
        Name = "web-server ${humid.server.id}"
    }

    # Read the AMI id "through" the humid resource to ensure that
    # both will change together.
    ami = humid.server.keepers.ami_id

    # ... (other aws_instance arguments) ...
}