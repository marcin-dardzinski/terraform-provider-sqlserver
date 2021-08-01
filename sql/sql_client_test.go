package sql

import (
	"os"
	"testing"
)

const tenantIdKey = "TFSQL_AZURE_TENANT_ID"
const clientIdKey = "TFSQL_AZURE_CLIENT_ID"
const clientSecretKey = "TFSQL_AZURE_CLIENT_SECRET"

func TestAuthorizesViaClientSecret(t *testing.T) {
	connString := &ConnectionString{
		ServerAddress: "tf-2137.database.windows.net",
		Database:      "tf-2137",
	}

	azure := &AzureADConfig{
		ClientId:     getEnv(clientIdKey, t),
		ClientSecret: getEnv(clientSecretKey, t),
		TenantId:     getEnv(tenantIdKey, t),
	}

	db, err := createUsingAzureActiveDirectoryAuth(connString, azure)

	if err != nil {
		t.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
}

func getEnv(name string, t *testing.T) string {
	val := os.Getenv(name)

	if val == "" {
		t.Fatalf("Error, missing environment variable %s", name)
	}
	return val
}
