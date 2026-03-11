package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingsSlack_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "akmatori_settings_slack" "test" {
  alerts_channel = "#test-channel"
  enabled        = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_settings_slack.test", "alerts_channel", "#test-channel"),
					resource.TestCheckResourceAttr("akmatori_settings_slack.test", "enabled", "false"),
				),
			},
		},
	})
}
