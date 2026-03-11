resource "akmatori_alert_source" "alertmanager" {
  source_type_name = "alertmanager"
  name             = "production-alertmanager"
  description      = "Production Prometheus Alertmanager"
  webhook_secret   = var.alertmanager_secret
  enabled          = true
}

variable "alertmanager_secret" {
  type      = string
  sensitive = true
}
