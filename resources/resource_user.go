package resources

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/marcin-dardzinski/terraform-provider-sqlserver/sql"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		Create:   resourceUserCreate,
		Read:     resourceUserRead,
		Update:   resourceUserUpdate,
		Delete:   resourceUserDelete,
		Importer: &schema.ResourceImporter{State: resourceUserImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"external"},
			},
			"external": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sql.SqlClient)
	userClient := sql.CreateSqlUserClient(client)

	name := d.Get("name").(string)
	password := d.Get("password").(string)

	err := userClient.Create(name, password)
	if err != nil {
		return err
	}
	d.SetId(client.Id + "/" + name)

	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sql.SqlClient)
	userClient := sql.CreateSqlUserClient(client)

	user, err := userClient.Get(d.Get("name").(string))
	if err != nil {
		return err
	}

	if user == nil {
		d.SetId("")
		return nil
	}

	if err = d.Set("name", user.Name); err != nil {
		return err
	}

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sql.SqlClient)
	userClient := sql.CreateSqlUserClient(client)

	name := d.Get("name").(string)

	d.Partial(true)

	if err := tryChangePassword(d, userClient, name); err != nil {
		return err
	}

	d.Partial(false)

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(sql.SqlUserClient)
	name := d.Get("name").(string)

	return client.Delete(name)
}

func resourceUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	client := m.(*sql.SqlClient)
	userClient := sql.CreateSqlUserClient(client)

	name := getUserNameFromId(d.Id())

	user, err := userClient.Get(name)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	if user == nil {
		d.SetId("")
		return []*schema.ResourceData{}, nil
	}

	d.SetId(client.Id + "/" + name)

	if err = d.Set("name", user.Name); err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}

func tryChangePassword(d *schema.ResourceData, client sql.SqlUserClient, name string) error {
	if d.HasChange("password") {
		_, new := d.GetChange("password")
		if err := client.ChangePassword(name, new.(string)); err != nil {
			return err
		}

		// TODO: Check if needed
		// d.SetPartial("password")
	}
	return nil
}

func getUserNameFromId(id string) string {
	s := strings.Split(id, "/")
	return s[len(s)-1]
}
