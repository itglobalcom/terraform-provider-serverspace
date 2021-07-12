variable "ubuntu" {
  description = "Default LTS"
  default     = "Ubuntu-18.04-X64"
}

variable "am_location" {
  description = "am location"
  default     = "am2"
}

variable "ssh_key_path" {
  type        = string
  description = "The file path to an ssh public key"
  default     = "~/.ssh/id_rsa.pub"
}


variable "ssh_private_key_path" {
  type        = string
  description = "The file path to an ssh privte key"
  default     = "~/.ssh/id_rsa"
}
