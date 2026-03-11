package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccToolInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "akmatori_tool_type" "test" {
  name = "zabbix"
}

resource "akmatori_tool_instance" "test" {
  tool_type_id = data.akmatori_tool_type.test.id
  name         = "tf-test-instance"
  enabled      = true
  settings_json = jsonencode({
    url      = "https://zabbix.test.com"
    username = "test"
    password = "test"
  })
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_tool_instance.test", "name", "tf-test-instance"),
					resource.TestCheckResourceAttr("akmatori_tool_instance.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("akmatori_tool_instance.test", "id"),
				),
			},
			{
				ResourceName:            "akmatori_tool_instance.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings_json"},
			},
		},
	})
}
