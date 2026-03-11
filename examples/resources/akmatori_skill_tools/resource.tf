resource "akmatori_skill_tools" "example" {
  skill_name        = akmatori_skill.example.name
  tool_instance_ids = [akmatori_tool_instance.example.id]
}
