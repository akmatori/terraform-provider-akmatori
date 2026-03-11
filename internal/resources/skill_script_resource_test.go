package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSkillScript_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "akmatori_skill" "script_test" {
  name    = "tf-test-skill-script"
  enabled = true
}

resource "akmatori_skill_script" "test" {
  skill_name = akmatori_skill.script_test.name
  filename   = "test.py"
  content    = "print('hello')"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_skill_script.test", "filename", "test.py"),
					resource.TestCheckResourceAttr("akmatori_skill_script.test", "content", "print('hello')"),
				),
			},
		},
	})
}
