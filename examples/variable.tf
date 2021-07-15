variable "ubuntu" {
  type        = string
  description = "Default LTS"
  default     = "Ubuntu-18.04-X64"
}

variable "am_location" {
  type        = string
  description = "am location"
  default     = "am2"
}

variable "kz_location" {
  type        = string
  description = "kz location"
  default     = "kz"
}

variable "ssh_key_path" {
  type        = string
  description = "The file path to an ssh public key"
  default     = "~/.ssh/id_rsa.pub"
}


variable "ssh_private_key_path" {
  type        = string
  sensitive   = true
  description = "The file path to an ssh privte key"
  default     = "~/.ssh/id_rsa"
}
