package resources_test

import (
	"fmt"
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

func TestAccSkill_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSkillConfig("tf-test-skill", "Test skill", "testing"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_skill.test", "name", "tf-test-skill"),
					resource.TestCheckResourceAttr("akmatori_skill.test", "description", "Test skill"),
					resource.TestCheckResourceAttr("akmatori_skill.test", "category", "testing"),
					resource.TestCheckResourceAttr("akmatori_skill.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("akmatori_skill.test", "id"),
				),
			},
			{
				ResourceName:      "akmatori_skill.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSkillConfig("tf-test-skill", "Updated skill", "updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_skill.test", "description", "Updated skill"),
					resource.TestCheckResourceAttr("akmatori_skill.test", "category", "updated"),
				),
			},
		},
	})
}

func testAccSkillConfig(name, description, category string) string {
	return fmt.Sprintf(`
resource "akmatori_skill" "test" {
  name        = %[1]q
  description = %[2]q
  category    = %[3]q
  enabled     = true
  prompt      = "Test prompt"
}
`, name, description, category)
}
