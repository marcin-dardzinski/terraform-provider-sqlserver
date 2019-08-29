package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func configSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_string": &schema.Schema{
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"database_id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the database to be prepended to resources' ids.\nIf not set <server_address>/<database_name> will be used",
		},
	}
}
