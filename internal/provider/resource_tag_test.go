package provider_test

import (
	"fmt"
	"testing"

	"github.com/gumieri/xyops-integrations/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTagConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("xyops_tag.test", "title", "test-tag"),
					resource.TestCheckResourceAttrSet("xyops_tag.test", "id"),
				),
			},
			{
				Config: testAccTagConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("xyops_tag.test", "title", "test-tag-updated"),
				),
			},
			{
				ResourceName:      "xyops_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}
}

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"xyops": providerserver.NewProtocol6WithError(provider.New("test")),
	}
}

func testAccImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource Not Found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

const testAccTagConfigBasic = `
provider "xyops" {
  base_url = "http://localhost:8000"
  api_key  = "test-api-key"
}

resource "xyops_tag" "test" {
  title = "test-tag"
}
`

const testAccTagConfigUpdate = `
provider "xyops" {
  base_url = "http://localhost:8000"
  api_key  = "test-api-key"
}

resource "xyops_tag" "test" {
  title = "test-tag-updated"
}
`
