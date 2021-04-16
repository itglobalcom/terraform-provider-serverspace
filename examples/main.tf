terraform {
  required_providers {
    serverspace = {
      version = "0.2"
      source  = "serverspace.by/main/serverspace"
    }
  }
}


variable "ubuntu" {
  description = "Default LTS"
  default     = "Ubuntu-18.04-X64"
}

variable "am_loc" {
  description = "am location"
  default     = "am2"
}


resource "serverspace_isolated_network" "my_net" {
  location    = var.am_loc
  name        = "my_net"
  description = "Internal network"
}

resource "serverspace_server" "vm1" {
  image    = var.ubuntu
  name     = "vm1"
  location = var.am_loc
  cpu      = 1
  ram      = 4

  volume "root" { # The name of the first volume block is ignored (like in CLI)
    size = 25
  }
  volume "bar" {
    size = 250
  }

  nic {
    bandwidth = 50
    type      = "shared"
  }
  nic {
    network = data.serverspace_isolated_network.my_net.id
    type    = "isolated"
  }

  ssh_keys = [
    data.serverspace_ssh_key.terraform.id
  ]

  tags = [
    "nginx", "proxy"
  ]

  connection {
    host        = self.public_ip_addresses[0] # Read-only attribute computed from connected networks
    user        = "root"
    type        = "ssh"
    private_key = file(var.pvt_key)
    timeout     = "2m"
  }
  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      # install nginx
      "sudo apt-get update",
      "sudo apt-get -y install nginx"
    ]
  }
}
