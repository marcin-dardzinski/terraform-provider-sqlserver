package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataConnectionString() *schema.Resource {
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

	conn := ConnectionString{
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
	id := fmt.Sprintf("%s/%d/%s/%s", server, port, database, username)

	d.SetId(id)
	d.Set("value", conn.String())

	return nil
}
