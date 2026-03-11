package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSkillTools_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "akmatori_skill" "tools_test" {
  name    = "tf-test-skill-tools"
  enabled = true
}

data "akmatori_tool_type" "tools_test" {
  name = "zabbix"
}

resource "akmatori_tool_instance" "tools_test" {
  tool_type_id = data.akmatori_tool_type.tools_test.id
  name         = "tf-test-tools-instance"
  enabled      = true
  settings_json = jsonencode({
    url      = "https://zabbix.test.com"
    username = "test"
    password = "test"
  })
}

resource "akmatori_skill_tools" "test" {
  skill_name        = akmatori_skill.tools_test.name
  tool_instance_ids = [akmatori_tool_instance.tools_test.id]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_skill_tools.test", "skill_name", "tf-test-skill-tools"),
				),
			},
		},
	})
}
