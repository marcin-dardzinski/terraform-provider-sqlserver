package sql

import (
	"database/sql"
	"os"
	"testing"

	"github.com/marcin-dardzinski/terraform-provider-sqlserver/consts"
)

func TestAuthorizeViaSQLUser(t *testing.T) {
	connString := tryGetEnv(consts.ConnectionStringEnv, t)
	db, err := createUsingPasswordAuth(connString)

	if err != nil {
		t.Fatal(err)
	}

	pingDb(db, t)
}

func TestAuthorizeViaClientSecret(t *testing.T) {
	connString := tryGetEnv(consts.UserlessConnectionStringEnv, t)
	azure := &AzureADConfig{
		ClientId:     tryGetEnv(consts.ClientIdEnv, t),
		ClientSecret: tryGetEnv(consts.ClientSecretEnv, t),
		TenantId:     tryGetEnv(consts.TenantIdEnv, t),
	}

	db, err := createUsingAzureActiveDirectoryAuth(connString, azure)

	if err != nil {
		t.Fatal(err)
	}

	pingDb(db, t)
}

func TestAuthorizeViaCli(t *testing.T) {
	connString := tryGetEnv(consts.UserlessConnectionStringEnv, t)
	azure := &AzureADConfig{
		SubscriptionId: os.Getenv(consts.SubscriptionIdEnv),
		UseCLI:         true,
	}

	db, err := createUsingAzureActiveDirectoryAuth(connString, azure)

	if err != nil {
		t.Fatal(err)
	}

	pingDb(db, t)
}

func pingDb(db *sql.DB, t *testing.T) {
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
}

func tryGetEnv(name string, t *testing.T) string {
	val, ok := os.LookupEnv(name)

	if !ok {
		t.Skipf("Missing environment variable %s, skipping test", name)
	}
	return val
}
