resource "akmatori_context_file" "runbook" {
  filename    = "runbook.md"
  content     = file("${path.module}/files/runbook.md")
  description = "Incident response runbook"
}
