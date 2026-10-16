package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gumieri/xyops-integrations/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTagConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("xyops_tag.test", "title", "test-tag"),
					resource.TestCheckResourceAttrSet("xyops_tag.test", "id"),
				),
			},
			{
				Config: testAccTagConfigUpdate(),
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

	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set; skipping acceptance test")
	}

	if os.Getenv("XYOPS_BASE_URL") == "" || os.Getenv("XYOPS_API_KEY") == "" {
		t.Fatal("XYOPS_BASE_URL and XYOPS_API_KEY must be set for acceptance tests")
	}
}

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"xyops": providerserver.NewProtocol6WithError(provider.New("test")),
	}
}

// testAccProviderConfig renders the provider block from the live-instance
// env vars (validated by testAccPreCheck).
func testAccProviderConfig() string {
	return fmt.Sprintf(`
provider "xyops" {
  base_url = %q
  api_key  = %q
}
`, os.Getenv("XYOPS_BASE_URL"), os.Getenv("XYOPS_API_KEY"))
}

func testAccTagConfigBasic() string {
	return testAccProviderConfig() + `
resource "xyops_tag" "test" {
  title = "test-tag"
}
`
}

func testAccTagConfigUpdate() string {
	return testAccProviderConfig() + `
resource "xyops_tag" "test" {
  title = "test-tag-updated"
}
`
}
