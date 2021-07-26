package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataConnectionString() *schema.Resource {
	return &schema.Resource{
		Read: dataConnectionStringRead,
		Schema: map[string]*schema.Schema{
			"server": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1433,
			},
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"connection_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"max_pool_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
			"trust_server_certificate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"persist_security_info": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"encrypt": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"value": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataConnectionStringRead(d *schema.ResourceData, m interface{}) error {
	conn := ExtractConnString(d)
	connString, err := conn.String()
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%s/%d/%s/%s", conn.ServerAddress, conn.Port, conn.Database, conn.Username)
	d.SetId(id)
	if err = d.Set("value", connString); err != nil {
		return err
	}

	return nil
}
