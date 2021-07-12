terraform {
  required_providers {
    serverspace = {
      version = "0.2"
      source  = "serverspace.by/main/serverspace"
    }
  }
}


resource "serverspace_ssh" "my_ssh" {
  name       = "just a key"
  public_key = file(var.ssh_key_path)
}

resource "serverspace_isolated_network" "my_net" {
  location       = var.am_location
  name           = "my_net"
  description    = "Internal network"
  network_prefix = "192.168.0.0"
  mask           = 24
}

resource "serverspace_server" "vm1" {
  name     = "vm1"
  image    = var.ubuntu
  location = var.am_location
  cpu      = 1
  ram      = 2048

  boot_volume_size = 30720 # 25600

  volume {
    name = "bar1"
    size = 30720
  }

  # public_nic {
  #   bandwidth = 50
  # }

  # private_nic {
  #   network = serverspace_isolated_network.my_net.id
  # }

  nic {
    # network = serverspace_isolated_network.my_net.id
    network_type = "PublicShared"
    bandwidth    = 50
  }

  connection {
    host        = self.public_nic[0].ip_address # Read-only attribute computed from connected networks
    user        = "root"
    type        = "ssh"
    private_key = file(var.ssh_private_key_path)
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

# output "vm1" {
#   value = serverspace_server.vm1
# }


# output "my_net" {
#   value = serverspace_isolated_network.my_net
# }
