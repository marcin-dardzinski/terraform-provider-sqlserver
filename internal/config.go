package internal

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func configSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_string": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			DefaultFunc: schema.EnvDefaultFunc(ConnectionStringEnv, nil),
		},
	}
}
