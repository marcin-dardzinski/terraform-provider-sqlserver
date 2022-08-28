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

	azureRaw := data.Get("azure").([]interface{})

	var azure *sql.AzureADConfig = nil

	if len(azureRaw) == 1 {
		azureRaw := azureRaw[0].(map[string]interface{})

		azure = &sql.AzureADConfig{
			TenantId:            azureRaw["tenant_id"].(string),
			SubscriptionId:      azureRaw["subscription_id"].(string),
			ClientId:            azureRaw["client_id"].(string),
			ClientSecret:        azureRaw["client_secret"].(string),
			CertificatePath:     azureRaw["client_certificate_path"].(string),
			CertificatePassword: azureRaw["client_certificate_password"].(string),
			UseMSI:              azureRaw["use_msi"].(bool),
			UseCLI:              azureRaw["use_cli"].(bool),
		}
	}

	return sql.CreatePooledSqlClient(sql.SqlClientConfig{
		ConnectionString: connString,
		Azure:            azure,
	})
}
