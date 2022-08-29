package internal

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/marcin-dardzinski/terraform-provider-sqlserver/internal/resources"
	"github.com/marcin-dardzinski/terraform-provider-sqlserver/sql"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"sqlserver_user":      resources.ResourceUser(),
			"sqlserver_user_role": resources.ResourceUserRole(),
		},
		Schema:        configSchema(),
		ConfigureFunc: createSqlClientProvider,
	}
}

func createSqlClientProvider(data *schema.ResourceData) (interface{}, error) {
	connString := data.Get("connection_string").(string)

	return sql.CreatePooledSqlClient(sql.SqlClientConfig{
		ConnectionString: connString,
	})
}
