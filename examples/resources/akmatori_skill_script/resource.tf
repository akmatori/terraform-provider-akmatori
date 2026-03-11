resource "akmatori_skill_script" "example" {
  skill_name = akmatori_skill.example.name
  filename   = "handler.py"
  content    = file("${path.module}/scripts/handler.py")
}
