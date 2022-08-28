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
			DefaultFunc: schema.EnvDefaultFunc(ConnectionStringEnv, ""),
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
						DefaultFunc: schema.EnvDefaultFunc(TenantIdEnv, ""),
					},
					"subscription_id": {
						Type:        schema.TypeString,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc(SubscriptionIdEnv, ""),
					},
					"client_id": {
						Type:        schema.TypeString,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc(ClientIdEnv, ""),
					},
					"client_secret": {
						Type:        schema.TypeString,
						Optional:    true,
						Sensitive:   true,
						DefaultFunc: schema.EnvDefaultFunc(ClientSecretEnv, ""),
					},
					"client_certificate_path": {
						Type:        schema.TypeString,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc(ClientCertificatePathEnv, ""),
					},
					"client_certificate_password": {
						Type:        schema.TypeString,
						Optional:    true,
						Sensitive:   true,
						DefaultFunc: schema.EnvDefaultFunc(ClientCertificatePasswordEnv, ""),
					},
					"use_msi": {
						Type:        schema.TypeBool,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc(UseMsiEnv, "false"),
					},
					"use_cli": {
						Type:        schema.TypeBool,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc(UseCliEnv, "true"),
					},
				},
			},
		},
	}
}

// func boolFromString(x string) bool {
// 	x = strings.ToLower(x)
// 	return x == "true" || x == "1" || x == "t" || x == "y"
// }
