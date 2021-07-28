package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/marcin-dardzinski/terraform-provider-sqlserver/resources"
	"github.com/marcin-dardzinski/terraform-provider-sqlserver/sql"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"sqlserver_user": resources.ResourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sqlserver_connection_string": resources.DataConnectionString(),
		},
		Schema:        configSchema(),
		ConfigureFunc: createSqlClientProvider,
	}
}

func createSqlClientProvider(data *schema.ResourceData) (interface{}, error) {
	connString := resources.ExtractConnString(data)
	return sql.CreateSqlClient(&connString)
}
