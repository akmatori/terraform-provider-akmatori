resource "akmatori_settings_aggregation" "main" {
  enabled                          = true
  correlation_confidence_threshold = 0.70
  merge_confidence_threshold       = 0.75
  recorrelation_enabled            = true
  recorrelation_interval_minutes   = 3
  max_incidents_to_analyze         = 20
  observing_duration_minutes       = 30
  correlator_timeout_seconds       = 5
  merge_analyzer_timeout_seconds   = 30
}
