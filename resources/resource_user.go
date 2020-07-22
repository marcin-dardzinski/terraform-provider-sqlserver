package resources

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/marcin-dardzinski/sql-server-terraform/sql"
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
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"roles": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(sql.SqlUserClient)

	name := d.Get("name").(string)
	password := d.Get("password").(string)
	roles := d.Get("roles").(*schema.Set)

	err := client.Create(name, password, castStrings(roles))
	if err != nil {
		return err
	}
	d.SetId(client.DatabaseId() + "/" + name)

	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(sql.SqlUserClient)

	user, err := client.Get(d.Get("name").(string))
	if err != nil {
		return err
	}

	if user == nil {
		d.SetId("")
		return nil
	}

	desiredRoles := d.Get("roles").(*schema.Set)
	roles := schema.NewSet(desiredRoles.F, []interface{}{})
	for _, r := range user.Roles {
		roles.Add(r)
	}
	knownRoles := desiredRoles.Intersection(roles)

	if err = d.Set("name", user.Name); err != nil {
		return err
	}

	if err = d.Set("roles", knownRoles); err != nil {
		return err
	}

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(sql.SqlUserClient)
	name := d.Get("name").(string)

	d.Partial(true)

	if err := tryChangePassword(d, client, name); err != nil {
		return err
	}

	if err := tryChangeRoles(d, client, name); err != nil {
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
	client := m.(sql.SqlUserClient)
	name := getUserNameFromId(d.Id())

	user, err := client.Get(name)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	if user == nil {
		d.SetId("")
		return []*schema.ResourceData{}, nil
	}

	d.SetId(client.DatabaseId() + "/" + name)

	if err = d.Set("name", user.Name); err != nil {
		return []*schema.ResourceData{}, err
	}

	roles := schema.NewSet(schema.HashString, []interface{}{})
	for _, role := range user.Roles {
		roles.Add(role)
	}

	if err = d.Set("roles", roles); err != nil {
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

		d.SetPartial("password")
	}
	return nil
}

func tryChangeRoles(d *schema.ResourceData, client sql.SqlUserClient, name string) error {
	if d.HasChange("roles") {
		oldRaw, newRaw := d.GetChange("roles")
		old, new := oldRaw.(*schema.Set), newRaw.(*schema.Set)

		grant := new.Difference(old)
		revoke := old.Difference(new)

		if err := client.ChangeRoles(name, castStrings(grant), castStrings(revoke)); err != nil {
			return err
		}

		d.SetPartial("roles")
	}

	return nil
}

func castStrings(set *schema.Set) []string {
	raw := set.List()
	result := make([]string, set.Len())
	for i := range raw {
		result[i] = raw[i].(string)
	}

	return result
}

func getUserNameFromId(id string) string {
	s := strings.Split(id, "/")
	return s[len(s)-1]
}
