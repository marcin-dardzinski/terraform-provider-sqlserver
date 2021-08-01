package main

import (
	"fmt"

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
	connStringRaw, ok := data.GetOk("connection_string")

	if !ok {
		return nil, fmt.Errorf("no connection string")
	}

	connString, err := sql.ParseConnectionString(connStringRaw.(string))

	if err != nil {
		return nil, err
	}

	azureRaw := data.Get("azure")

	return sql.CreatePooledSqlClient(sql.SqlClientConfig{
		ConnectionString: connString,
		Azure:            azureRaw.(*sql.AzureADConfig),
	})
}
