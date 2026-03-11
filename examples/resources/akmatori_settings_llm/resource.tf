resource "akmatori_settings_llm" "anthropic" {
  llm_provider   = "anthropic"
  api_key        = var.anthropic_api_key
  model          = "claude-sonnet-4-20250514"
  thinking_level = "medium"
}

variable "anthropic_api_key" {
  type      = string
  sensitive = true
}
