package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingsAggregation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "akmatori_settings_aggregation" "test" {
  enabled                          = true
  correlation_confidence_threshold = 0.80
  merge_confidence_threshold       = 0.85
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_settings_aggregation.test", "enabled", "true"),
					resource.TestCheckResourceAttr("akmatori_settings_aggregation.test", "correlation_confidence_threshold", "0.8"),
				),
			},
		},
	})
}
