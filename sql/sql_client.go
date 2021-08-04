package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/denisenkom/go-mssqldb/msdsn"
)

type SqlClientConfig struct {
	ConnectionString string
	Azure            *AzureADConfig
}

type AzureADConfig struct {
	TenantId            string
	SubscriptionId      string
	ClientId            string
	ClientSecret        string
	CertificatePath     string
	CertificatePassword string
	UseMSI              bool
	UseCLI              bool
}

type SqlClient struct {
	Db *sql.DB
	Id string
}

func CreateSqlClient(config SqlClientConfig) (*SqlClient, error) {
	var db *sql.DB

	parsedConnStr, _, err := msdsn.Parse(config.ConnectionString)

	if err != nil {
		return nil, err
	}

	if config.Azure == nil {
		db, err = createUsingPasswordAuth(config.ConnectionString)
	} else {
		db, err = createUsingAzureActiveDirectoryAuth(config.ConnectionString, config.Azure)
	}

	if err != nil {
		return nil, err
	}

	return &SqlClient{
		Db: db,
		Id: parseDatabaseId(&parsedConnStr),
	}, nil
}

func parseDatabaseId(config *msdsn.Config) string {
	return config.Host + "/" + config.Database
}

func createUsingPasswordAuth(connString string) (*sql.DB, error) {
	return sql.Open("mssql", connString)
}

func createUsingAzureActiveDirectoryAuth(connString string, azure *AzureADConfig) (*sql.DB, error) {
	var cred azcore.TokenCredential
	var err error

	if azure.ClientSecret != "" {
		cred, err = azidentity.NewClientSecretCredential(azure.TenantId, azure.ClientId, azure.ClientSecret, nil)
	} else if azure.CertificatePath != "" {
		cred, err = azidentity.NewClientCertificateCredential(azure.TenantId, azure.ClientId, azure.CertificatePath, &azidentity.ClientCertificateCredentialOptions{
			Password: azure.CertificatePassword,
		})
	} else if azure.UseMSI {
		cred, err = azidentity.NewManagedIdentityCredential(azure.ClientId, nil)
	} else if azure.UseCLI {
		cred, err = azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{
			TokenProvider: tenantAwareAzureCLITokenProvider(azure.TenantId, azure.SubscriptionId),
		})
	} else {
		err = errors.New("no azure authentication method selected")
	}

	if err != nil {
		return nil, err
	}

	connector, err := mssql.NewAccessTokenConnector(connString, func() (string, error) {
		token, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{
			Scopes: []string{"https://database.windows.net/.default"},
		})

		if err != nil {
			return "", err
		}

		return token.Token, nil
	})

	if err != nil {
		return nil, err
	}

	return sql.OpenDB(connector), nil
}
