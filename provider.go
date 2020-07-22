package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"sqlserver_user": resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sqlserver_connection_string": dataConnectionString(),
		},
		Schema:        configSchema(),
		ConfigureFunc: createSqlClientProvider,
	}
}

func createSqlClientProvider(data *schema.ResourceData) (interface{}, error) {
	connString := data.Get("connection_string").(string)

	if connString == "" {
		x := extractConnString(data)
		var err error
		connString, err = x.String()

		if err != nil {
			return nil, err
		}
	}

	return CreateSqlClient(connString)
}
