package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSHKey_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "akmatori_tool_type" "ssh_test" {
  name = "ssh"
}

resource "akmatori_tool_instance" "ssh_test" {
  tool_type_id = data.akmatori_tool_type.ssh_test.id
  name         = "tf-test-ssh-instance"
  enabled      = true
}

resource "akmatori_ssh_key" "test" {
  tool_id     = akmatori_tool_instance.ssh_test.id
  name        = "tf-test-key"
  private_key = "-----BEGIN OPENSSH PRIVATE KEY-----\ntest\n-----END OPENSSH PRIVATE KEY-----"
  is_default  = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("akmatori_ssh_key.test", "name", "tf-test-key"),
					resource.TestCheckResourceAttrSet("akmatori_ssh_key.test", "key_id"),
				),
			},
		},
	})
}
