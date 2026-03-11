resource "akmatori_settings_proxy" "main" {
  proxy_url      = "http://proxy.example.com:8080"
  no_proxy       = "localhost,127.0.0.1"
  openai_enabled = true
  slack_enabled  = true
  zabbix_enabled = false
}
