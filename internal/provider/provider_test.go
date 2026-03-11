package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"akmatori": providerserver.NewProtocol6WithError(New("test")()),
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
		if p := os.Getenv("AKMATORI_PASSWORD"); p == "" {
			t.Fatal("AKMATORI_PASSWORD must be set when using AKMATORI_USERNAME")
		}
	}
}
