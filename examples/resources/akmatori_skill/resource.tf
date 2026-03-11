resource "akmatori_skill" "zabbix_analyst" {
  name        = "zabbix-analyst"
  description = "Analyzes Zabbix alerts and provides diagnosis"
  category    = "monitoring"
  enabled     = true

  prompt = <<-EOT
    You are a Zabbix monitoring analyst.
    Analyze the incoming alert and provide a diagnosis.
  EOT
}

resource "akmatori_skill_tools" "zabbix_analyst_tools" {
  skill_name       = akmatori_skill.zabbix_analyst.name
  tool_instance_ids = [akmatori_tool_instance.zabbix.id]
}
