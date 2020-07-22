package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func configSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_string": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
	}
}
