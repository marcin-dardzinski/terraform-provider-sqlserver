package main

import "github.com/hashicorp/terraform/helper/schema"

func extractConnString(d *schema.ResourceData) ConnectionString {
	server := d.Get("server").(string)
	port := d.Get("port").(int)
	database := d.Get("database").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	connectionTimeout := d.Get("connection_timeout").(int)
	maxPoolSize := d.Get("max_pool_size").(int)
	trustServerCertificate := d.Get("trust_server_certificate").(bool)
	persistSecurityInfo := d.Get("persist_security_info").(bool)
	encrypt := d.Get("encrypt").(bool)

	return ConnectionString{
		ServerAddress:          server,
		Port:                   port,
		Database:               database,
		Username:               username,
		Password:               password,
		ConnectionTimeout:      connectionTimeout,
		MaxPoolSize:            maxPoolSize,
		TrustServerCertificate: trustServerCertificate,
		PersistSecurityInfo:    persistSecurityInfo,
		Encrypt:                encrypt,
	}
}
