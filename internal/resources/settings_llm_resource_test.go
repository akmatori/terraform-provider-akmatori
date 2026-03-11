package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettingsLLM_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "akmatori_settings_llm" "test" {
  llm_provider   = "anthropic"
  api_key        = "test-key"
  model          = "claude-sonnet-4-20250514"
  thinking_level = "medium"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_settings_llm.test", "llm_provider", "anthropic"),
					resource.TestCheckResourceAttr("akmatori_settings_llm.test", "model", "claude-sonnet-4-20250514"),
				),
			},
		},
	})
}
