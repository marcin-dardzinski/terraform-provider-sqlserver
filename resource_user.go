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
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(SqlUserClient)

	name := d.Get("name").(string)
	password := d.Get("password").(string)

	err := client.Create(name, password)
	if err != nil {
		return err
	}
	d.SetId(name)

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

	d.Set("name", user.name)

	return nil

}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(SqlUserClient)
	name := d.Get("name").(string)

	d.Partial(true)

	if err := tryChangePassword(d, client, name); err != nil {
		return nil
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
