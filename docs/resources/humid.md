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

## Schema

### Optional

- **keepers** (Map of String) Arbitrary map of values that, when changed, will trigger recreation of resource. See [the main provider documentation](../index.html) for more information.
- **wordlist** (String) Name of the wordlist to use to generate the random ID. See the [humid repo](https://github.com/kscarlett/humid/tree/main/wordlist) for more information on the wordlists available. Defaults to the "animals" list.
- **adjectives** (Integer) Amount of adjectives to use to generate the ID. Adds a lot more options. Defaults to 1.
- **separator** (String) What to use between words in the ID. Defaults to '-'.
- **capitalize** (Boolean) Whether to capitalize the first letter of each word or to leave everything lowercase. Defaults to false (lowercase).

### Read-Only

- **id** (String) The generated id presented in string format.
- **result** (String) The generated id presented in string format.

## Import

Import is supported using the following syntax:

```shell
# Random UUID's can be imported. This can be used to replace a config
# value with a value interpolated from the random provider without
# experiencing diffs.

terraform import humid.main cloudy-stylish-barb
```