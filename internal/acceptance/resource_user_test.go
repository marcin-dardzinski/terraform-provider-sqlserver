package resources

import (
	"fmt"
	"testing"

	"github.com/marcin-dardzinski/terraform-provider-sqlserver/internal"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccUser_basic(t *testing.T) {
	username := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	password1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) + "a1!"
	password2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) + "a1!"

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sqlserver": func() (*schema.Provider, error) { return internal.Provider(), nil },
		},

		Steps: []resource.TestStep{
			{
				Config: resourceUserWithPassword(username, password1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sqlserver_user.user", "name", username),
					resource.TestCheckResourceAttr("sqlserver_user.user", "password", password1),
				),
			},
			{
				Config: resourceUserWithPassword(username, password2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sqlserver_user.user", "name", username),
					resource.TestCheckResourceAttr("sqlserver_user.user", "password", password1),
				),
			},
		},
	})
}

func resourceUserWithPassword(username, password string) string {
	template := `
	resource "sqlserver_user" "user" {
		name = "%s"
		password = "%s"
	}`

	return fmt.Sprintf(template, username, password)
}
