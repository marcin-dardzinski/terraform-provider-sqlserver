package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"sqlserver_user": resourceUser(),
		},
		Schema:        configSchema(),
		ConfigureFunc: createSqlClientProvider,
	}
}

func createSqlClientProvider(data *schema.ResourceData) (interface{}, error) {
	connString := data.Get("connection_string").(string)
	dbId := ""
	if tmp := data.Get("database_id"); tmp != nil {
		dbId = tmp.(string)
	}

	return GetSqlClient(connString, dbId)
}
