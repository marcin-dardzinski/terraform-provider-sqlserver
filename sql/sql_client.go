package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	mssql "github.com/denisenkom/go-mssqldb"
)

type SqlClientConfig struct {
	ConnectionString *ConnectionString
	Azure            *AzureADConfig
}

type AzureADConfig struct {
	TenantId            string
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
	var err error

	if config.Azure == nil {
		db, err = createUsingPasswordAuth(config.ConnectionString)
	} else {
		db, err = createUsingAzureActiveDirectoryAuth(config.ConnectionString, config.Azure)
	}

	if err != nil {
		return nil, err
	}

	id := parseDatabaseId(config.ConnectionString)

	return &SqlClient{
		Db: db,
		Id: id,
	}, nil
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

	var cred azcore.TokenCredential

	if azure.ClientSecret != "" {
		cred, err = azidentity.NewClientSecretCredential(azure.TenantId, azure.ClientId, azure.ClientSecret, nil)
	} else if azure.CertificatePath != "" {
		cred, err = azidentity.NewClientCertificateCredential(azure.TenantId, azure.ClientId, azure.CertificatePath, &azidentity.ClientCertificateCredentialOptions{
			Password: azure.CertificatePassword,
		})
	} else if azure.UseMSI {
		cred, err = azidentity.NewManagedIdentityCredential(azure.ClientId, nil)
	} else if azure.UseCLI {
		cred, err = azidentity.NewAzureCLICredential(nil)
	} else {
		err = errors.New("no azure authentication method selected")
	}

	if err != nil {
		return nil, err
	}

	connector, err := mssql.NewAccessTokenConnector(str, func() (string, error) {
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
