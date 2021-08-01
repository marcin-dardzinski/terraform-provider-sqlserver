package main

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func configSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_string": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			DefaultFunc: schema.EnvDefaultFunc("TFSQL_CONNECTION_STRING", ""),
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
						DefaultFunc: schema.EnvDefaultFunc("TFSQL_AZURE_TENANT_ID", ""),
					},
					"client_id": {
						Type:        schema.TypeString,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc("TFSQL_AZURE_CLIENT_ID", ""),
					},
					"client_secret": {
						Type:        schema.TypeString,
						Optional:    true,
						Sensitive:   true,
						DefaultFunc: schema.EnvDefaultFunc("TFSQL_AZURE_CLIENT_SECRET", ""),
					},
					"client_certificate_path": {
						Type:        schema.TypeString,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc("TFSQL_AZURE_CLIENT_CERTIFICATE_PATH", ""),
					},
					"client_certificate_password": {
						Type:        schema.TypeString,
						Optional:    true,
						Sensitive:   true,
						DefaultFunc: schema.EnvDefaultFunc("TFSQL_AZURE_CLIENT_CERTIFICATE_PASSWORD", ""),
					},
					"use_msi": {
						Type:        schema.TypeBool,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc("TFSQL_AZURE_USE_MSI", "false"),
						// func() (interface{}, error) {
						// 	str, err := schema.EnvDefaultFunc("TFSQL_AZURE_USE_MSI", "false")
						// 	return boolFromString(str), err
						// },
					},
					"use_cli": {
						Type:        schema.TypeBool,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc("TFSQL_AZURE_USE_CLI", "true"),
					},
				},
			},
		},
	}
}

func boolFromString(x string) bool {
	x = strings.ToLower(x)
	return x == "true" || x == "1" || x == "t" || x == "y"
}
