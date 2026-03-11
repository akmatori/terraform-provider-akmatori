package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlertSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "akmatori_alert_source" "test" {
  source_type_name = "alertmanager"
  name             = "tf-test-alertsource"
  description      = "Test alert source"
  enabled          = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_alert_source.test", "name", "tf-test-alertsource"),
					resource.TestCheckResourceAttrSet("akmatori_alert_source.test", "uuid"),
				),
			},
			{
				ResourceName:            "akmatori_alert_source.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"webhook_secret"},
			},
		},
	})
}
