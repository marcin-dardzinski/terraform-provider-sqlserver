package sql

import (
	"database/sql"
	"errors"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	mssql "github.com/denisenkom/go-mssqldb"
)

type SqlClientConfig struct {
	ConnectionString *ConnectionString
	Azure            *AzureADConfig
}

type AzureADConfig struct {
	TenantId     string
	ClientId     string
	ClientSecret string
}

func CreateDbConnection(config SqlClientConfig) (*sql.DB, error) {
	if config.Azure != nil {
		return createUsingPasswordAuth(config.ConnectionString)
	} else {
		return createUsingAzureActiveDirectoryAuth(config.ConnectionString, config.Azure)
	}
}

func createUsingPasswordAuth(connString *ConnectionString) (*sql.DB, error) {
	str, err := connString.String()
	if err != nil {
		return nil, err
	}

	return sql.Open("mssql", str)
}

func createUsingAzureActiveDirectoryAuth(connString *ConnectionString, azure *AzureADConfig) (*sql.DB, error) {
	if connString.Password != "" || connString.Username != "" {
		return nil, errors.New("connection string must not have username nor password when using Azure AD auth")
	}

	str, err := connString.String()
	if err != nil {
		return nil, err
	}

	config := auth.NewClientCredentialsConfig(azure.ClientId, azure.ClientSecret, azure.TenantId)
	config.Resource = "https://database.windows.net/"

	connector, err := mssql.NewAccessTokenConnector(str, func() (string, error) {
		token, err := config.ServicePrincipalToken()
		if err != nil {
			return "", err
		}
		if err = token.EnsureFresh(); err != nil {
			return "", err
		}

		return token.Token().AccessToken, nil
	})

	// connector, err := mssql.NewAccessTokenConnector(str, func() (string, error) {
	// 	token, err := cli.GetTokenFromCLI("https://database.windows.net/")
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	return token.AccessToken, nil
	// })

	if err != nil {
		return nil, err
	}

	return sql.OpenDB(connector), nil
}
