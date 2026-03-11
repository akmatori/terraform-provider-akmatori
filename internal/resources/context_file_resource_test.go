package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContextFile_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "akmatori_context_file" "test" {
  filename    = "tf-test-file.txt"
  content     = "test content"
  description = "Test context file"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_context_file.test", "filename", "tf-test-file.txt"),
					resource.TestCheckResourceAttrSet("akmatori_context_file.test", "id"),
				),
			},
		},
	})
}
