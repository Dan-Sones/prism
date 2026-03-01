variable "allowed_ip" {
  description = "Your public IP for SSH access"
  type        = string
}

variable "ssh_public_key" {
  description = "Path to SSH public key"
  type        = string
}
