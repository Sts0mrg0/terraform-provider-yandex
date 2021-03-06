---
layout: "yandex"
page_title: "Yandex: yandex_vpc_subnet"
sidebar_current: "docs-yandex-datasource-vpc-subnet"
description: |-
  Get information about a Yandex VPC subnet.
---

# yandex\_vpc\_subnet

Get information about a Yandex VPC subnet. See details about [Yandex Cloud VPC](https://cloud.yandex.com/docs/vpc/concepts/index).

```hcl
data "yandex_vpc_subnet" "admin" {
  subnet_id = "my-subnet-id"
}
```

This data source is used to define [VPC Subnets] to use to other resources.

## Argument Reference

The following arguments are supported:

* `subnet_id` (Required) - Subnet id.

## Attributes Reference

The following attribute is exported:

* `name` - Name of the subnet. 
* `description` - Description of the subnet.
* `folder_id` - ID of the folder in which the resource belongs.
* `network_id` - ID of the network this subnet belongs to.
* `labels` - Labels to assign to this subnet.
* `zone` - Name of the availability zone for this subnet.
* `v4_cidr_blocks` - The ranges of internal IPv4 addresses that are owned by this subnet.
* `v6_cidr_blocks` - The ranges of internal IPv6 addresses that are owned by this subnet.
* `created_at` - Creation timestamp of this subnet.

~> **Note:** `v6_cidr_blocks` attribute is currently not supported. It will be available in the future.

[VPC Subnets]: https://cloud.yandex.com/docs/vpc/concepts/network#subnet
