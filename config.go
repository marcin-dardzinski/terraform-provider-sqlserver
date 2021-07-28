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
			// ConflictsWith: []string{
			// 	"server",
			// 	"port",
			// 	"database",
			// 	"username",
			// 	"password",
			// 	"connection_timeout",
			// 	"max_pool_size",
			// 	"trust_server_certificate",
			// 	"persist_security_info",
			// 	"encrypt",
			// },
		},
		"server": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_SERVER", ""),
		},
		"port": {
			Type:        schema.TypeInt,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_SERVER", 1433),
		},
		"database": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_DATABASE", ""),
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_USERNAME", ""),
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_PASSWORD", ""),
		},
		"connection_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_CONNECTION_TIMEOUT", 30),
		},
		"max_pool_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_MAX_POOL_SIZE", 100),
		},
		"trust_server_certificate": {
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_TRUST_SERVER_CERTIFICATE", false),
		},
		"persist_security_info": {
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_PERSIST_SECURITY_INFO", false),
		},
		"encrypt": {
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SQLSERVER_ENCRYPT", true),
		},
		"azure": {
			Type:     schema.TypeMap,
			Optional: true,
		},
	}
}
