resource "akmatori_ssh_key" "example" {
  tool_id     = akmatori_tool_instance.example.id
  name        = "deploy-key"
  private_key = var.ssh_private_key
  is_default  = true
}

variable "ssh_private_key" {
  type      = string
  sensitive = true
}
