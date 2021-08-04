package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/marcin-dardzinski/terraform-provider-sqlserver/sql"
)

func ResourceUserRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserRoleCreate,
		Read:   resourceUserRoleRead,
		Delete: resourceUserRoleDelete,
		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceUserRoleCreate(d *schema.ResourceData, m interface{}) error {
	client, roleClient := makeUserRoleClient(m)

	role := getUserRoleModel(d)

	if err := roleClient.Grant(role); err != nil {
		return err
	}

	d.SetId(resourceUserRoleId(client.Id, role))
	return nil
}

func resourceUserRoleRead(d *schema.ResourceData, m interface{}) error {
	client, roleClient := makeUserRoleClient(m)

	role := getUserRoleModel(d)

	exists, err := roleClient.Exists(role)

	if err != nil {
		return err
	}

	if exists {
		d.SetId(resourceUserRoleId(client.Id, role))
	} else {
		d.SetId("")
	}

	return nil
}

func resourceUserRoleDelete(d *schema.ResourceData, m interface{}) error {
	_, roleClient := makeUserRoleClient(m)

	role := getUserRoleModel(d)

	return roleClient.Revoke(role)
}

func makeUserRoleClient(m interface{}) (*sql.SqlClient, sql.SqlUserRoleClient) {
	client := m.(*sql.SqlClient)
	roleClient := sql.CreateSqlUserRoleClient(client)
	return client, roleClient
}

func getUserRoleModel(d *schema.ResourceData) sql.SqlUserRole {
	user := d.Get("user").(string)
	role := d.Get("role").(string)

	return sql.SqlUserRole{
		User: user,
		Role: role,
	}
}

func resourceUserRoleId(dbId string, role sql.SqlUserRole) string {
	return dbId + "/user/" + role.User + "/role/" + role.Role
}
