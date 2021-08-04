package main

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/marcin-dardzinski/terraform-provider-sqlserver/consts"
)

func configSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_string": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			DefaultFunc: schema.EnvDefaultFunc(consts.ConnectionStringEnv, ""),
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
						DefaultFunc: schema.EnvDefaultFunc(consts.TenantIdEnv, ""),
					},
					"subscription_id": {
						Type:        schema.TypeString,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc(consts.SubscriptionIdEnv, ""),
					},
					"client_id": {
						Type:        schema.TypeString,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc(consts.ClientIdEnv, ""),
					},
					"client_secret": {
						Type:        schema.TypeString,
						Optional:    true,
						Sensitive:   true,
						DefaultFunc: schema.EnvDefaultFunc(consts.ClientSecretEnv, ""),
					},
					"client_certificate_path": {
						Type:        schema.TypeString,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc(consts.ClientCertificatePathEnv, ""),
					},
					"client_certificate_password": {
						Type:        schema.TypeString,
						Optional:    true,
						Sensitive:   true,
						DefaultFunc: schema.EnvDefaultFunc(consts.ClientCertificatePasswordEnv, ""),
					},
					"use_msi": {
						Type:        schema.TypeBool,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc(consts.UseMsiEnv, "false"),
					},
					"use_cli": {
						Type:        schema.TypeBool,
						Optional:    true,
						DefaultFunc: schema.EnvDefaultFunc(consts.UseCliEnv, "true"),
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
