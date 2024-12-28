package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "slash_shortcut" "test" {
  name  = "mb"
  link  = "https://foo.bar"
  title = "Microbin"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("slash_shortcut.test", "id"),
					resource.TestCheckResourceAttrSet("slash_shortcut.test", "name"),
					resource.TestCheckResourceAttrSet("slash_shortcut.test", "link"),
					resource.TestCheckResourceAttrSet("slash_shortcut.test", "title"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "slash_shortcut.test",
				ImportState:       true,
				ImportStateVerify: true,
				//// The last_updated attribute does not exist in the Slash
				//// API, therefore there is no value for it during import.
				//ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "slash_shortcut" "test" {
  name  = "mbs"
  link  = "https://foo.bar"
  title = "Microbin"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("slash_shortcut.test", "id"),
					resource.TestCheckResourceAttrSet("slash_shortcut.test", "name"),
					resource.TestCheckResourceAttrSet("slash_shortcut.test", "link"),
					resource.TestCheckResourceAttrSet("slash_shortcut.test", "title"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
