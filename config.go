package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func configSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_string": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_CONN_STRING", ""),
		},
		"azure": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"tenant_id": {
						Type:        schema.TypeString,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
					},
				},
			},
		},
	}
}
