---
page_title: "Netcup SCP Provider"
subcategory: ""
description: |-
  
---

# Provider for Netcup SCP API

This is a [Terraform](https://terraform.io) provider for the [Netcup](https://www.netcup.de/) SCP [webservice](https://www.netcup-wiki.de/wiki/Server_Control_Panel_(SCP)#Webservice).

## Example usage
```terraform
terraform {
  required_providers {
    netcup-ccp = {
      source = "rincedd/netcup-scp"
    }
  }
}

provider "netcup-scp" {
  login_name = "123456"     # Netcup customer number
  password   = "secret"     # SCP webservice password
}

data "netcup_vserver" "my_server" {
  server_name = "v12345678901234567"
}

```
