package sql

// import (
// 	"database/sql"
// 	"os"
// 	"testing"
// )

// func TestAuthorizeViaSQLUser(t *testing.T) {
// 	connString := tryGetEnv(ConnectionStringEnv, t)
// 	db, err := createUsingPasswordAuth(connString)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	pingDb(db, t)
// }

// func TestAuthorizeViaClientSecret(t *testing.T) {
// 	connString := tryGetEnv(UserlessConnectionStringEnv, t)
// 	azure := &AzureADConfig{
// 		ClientId:     tryGetEnv(ClientIdEnv, t),
// 		ClientSecret: tryGetEnv(ClientSecretEnv, t),
// 		TenantId:     tryGetEnv(TenantIdEnv, t),
// 	}

// 	db, err := createUsingAzureActiveDirectoryAuth(connString, azure)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	pingDb(db, t)
// }

// func TestAuthorizeViaCli(t *testing.T) {
// 	connString := tryGetEnv(UserlessConnectionStringEnv, t)
// 	azure := &AzureADConfig{
// 		SubscriptionId: os.Getenv(SubscriptionIdEnv),
// 		UseCLI:         true,
// 	}

// 	db, err := createUsingAzureActiveDirectoryAuth(connString, azure)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	pingDb(db, t)
// }

// func pingDb(db *sql.DB, t *testing.T) {
// 	if err := db.Ping(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func tryGetEnv(name string, t *testing.T) string {
// 	val, ok := os.LookupEnv(name)

// 	if !ok {
// 		t.Skipf("Missing environment variable %s, skipping test", name)
// 	}
// 	return val
// }
