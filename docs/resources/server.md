---
page_title: "serverspace_server Resource - terraform-provider-serverspace"
---

# serverspace_server (Resource)

Serverspace Server resource can be used to create, modify, and delete Servers. Servers also support [provisioning](https://www.terraform.io/docs/language/resources/provisioners/syntax.html).

## Example Usage

### Create a new server
```hcl
resource "serverspace_server" "example_server" {
  name = "example-server"
  image = "Ubuntu-20.04-X64"
  location = "am2"
  cpu = 1
  ram = 1024

  boot_volume_size = 25 * 1024
  
  volume {
    name = "additional-volume"
    size = 100 * 1024
  }
  
  nic {
    network = ""
    network_type = "PublicShared"
    bandwidth = 50
  }

  nic {
    network = resource.serverspace_isolated_network.example_net.id
    network_type = "Isolated"
    bandwidth = 0
  }
  
  ssh_keys = [
    resource.serverspace_ssh.example_key.id
  ]
```



## Schema

### Required

- **name** (String) Name of the server.
- **location** (String) Geographical location of the server.
- **image** (String) Image name (e.g. "Ubuntu-18.04-X64"), you can obtain it from API or CLI.
- **cpu** (Number) Count of the CPU Cores.
- **ram** (Number) Size of RAM in MB.
- **boot_volume_size** (Number) Size of the volume from which the server will be booted.
- **nic** (Block Set, Min: 1, Max: 5) (see [below for nested schema](#nestedblock--nic)) Network interface.

### Optional
- **volume** (Block List) (see [below for nested schema](#nestedblock--volume)) Additional volume description,
- **id** (String) The ID of this resource.
- **ssh_keys** (List of Number) List of IDs of ssh-keys.

### Read-Only

- **boot_volume_id** (Number) Id of boot volume.
- **public_ip_addresses** (List of String) List of assigned public IPs.

<a id="nestedblock--nic"></a>
### Nested Schema for `nic`

Required:

- **bandwidth** (Number) Network bandwidth in MB (for isolated networks should be 0).
- **network** (String) Network ID (for isolated networks).
- **network_type** (String) `PublicShared` or `Isolated`.

Read-Only:

- **id** (Number) The ID of this network interface.
- **ip_address** (String) Assigned address to the interface.


<a id="nestedblock--volume"></a>
### Nested Schema for `volume`

Required:

- **name** (String) Name of the volume.
- **size** (Number) Size of the volume in MB.

Read-Only:

- **id** (Number) The ID of this volume.




## Metadata information

For more details about the list of available locations and images, use the [s2ctl](https://github.com/itglobalcom/s2ctl) tool.

### Usage of s2ctl

- **locations** Returns a list of locations where server and network creation is available.
```
>s2ctl locations
- am2
- ds1
- nj3
- kz
```
- **images** Returns a list of available OS images which can be used to create a server.
```
>s2ctl images
- Freebsd-12.2-X64
- Debian-10.7-X64
- Windows-Server 2019 Std-X64
- Oracle-8.3-X64
- Ubuntu-18.04-X64
- CentOS-8.2-X64
- CentOS-7.9-X64
- Freebsd-13.0-X64
- Ubuntu-20.04-X64
- CentOS-8.3-X64
```
