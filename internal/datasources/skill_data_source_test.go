package datasources_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/akmatori/terraform-provider-akmatori/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"akmatori": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if v := os.Getenv("AKMATORI_HOST"); v == "" {
		t.Fatal("AKMATORI_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("AKMATORI_TOKEN"); v == "" {
		if u := os.Getenv("AKMATORI_USERNAME"); u == "" {
			t.Fatal("AKMATORI_TOKEN or AKMATORI_USERNAME must be set for acceptance tests")
		}
	}
}

func TestAccSkillDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "akmatori_skill" "ds_test" {
  name        = "tf-test-skill-ds"
  description = "Data source test skill"
  category    = "testing"
  enabled     = true
}

data "akmatori_skill" "test" {
  name = akmatori_skill.ds_test.name
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.akmatori_skill.test", "name", "tf-test-skill-ds"),
					resource.TestCheckResourceAttr("data.akmatori_skill.test", "description", "Data source test skill"),
				),
			},
		},
	})
}
