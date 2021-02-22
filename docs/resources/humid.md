---
page_title: "humid Resource - terraform-provider-humid"
subcategory: ""
description: |-
  The resource humid generates random human-friendly id strings that are intended to be used as unique identifiers for other resources.
---

# humid (Resource)

The resource humid generates random human-friendly id strings that are intended to be used as unique identifiers for other resources.

This resource uses [kscarlett/humid](https://github.com/kscarlett/humid) to generate unique names.

## Example usage

```terraform
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
```