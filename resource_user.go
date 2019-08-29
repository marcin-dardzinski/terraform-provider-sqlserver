package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"roles": &schema.Schema{
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
	client := m.(SqlUserClient)

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
	client := m.(SqlUserClient)

	user, err := client.Get(d.Get("name").(string))
	if err != nil {
		return err
	}

	if user == nil {
		d.SetId("")
		return nil
	}

	roles := schema.NewSet(schema.HashString, []interface{}{})
	for _, role := range user.roles {
		roles.Add(role)
	}

	d.Set("name", user.name)
	d.Set("roles", &roles)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(SqlUserClient)
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
	client := m.(SqlUserClient)
	name := d.Get("name").(string)

	return client.Delete(name)
}

func tryChangePassword(d *schema.ResourceData, client SqlUserClient, name string) error {
	if d.HasChange("password") {
		_, new := d.GetChange("password")
		if err := client.ChangePassword(name, new.(string)); err != nil {
			return err
		}

		d.SetPartial("password")
	}
	return nil
}

func tryChangeRoles(d *schema.ResourceData, client SqlUserClient, name string) error {
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
