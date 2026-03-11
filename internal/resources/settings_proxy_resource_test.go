package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingsProxy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "akmatori_settings_proxy" "test" {
  proxy_url      = "http://proxy.test:8080"
  no_proxy       = "localhost"
  openai_enabled = true
  slack_enabled  = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_settings_proxy.test", "proxy_url", "http://proxy.test:8080"),
					resource.TestCheckResourceAttr("akmatori_settings_proxy.test", "slack_enabled", "false"),
				),
			},
		},
	})
}
