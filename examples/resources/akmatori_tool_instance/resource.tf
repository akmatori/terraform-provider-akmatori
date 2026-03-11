data "akmatori_tool_type" "zabbix" {
  name = "zabbix"
}

resource "akmatori_tool_instance" "zabbix" {
  tool_type_id = data.akmatori_tool_type.zabbix.id
  name         = "production-zabbix"
  enabled      = true

  settings_json = jsonencode({
    url      = "https://zabbix.example.com"
    username = "api-user"
    password = var.zabbix_password
  })
}

variable "zabbix_password" {
  type      = string
  sensitive = true
}
