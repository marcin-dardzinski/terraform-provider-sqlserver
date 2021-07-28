package sql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure/cli"
	mssql "github.com/denisenkom/go-mssqldb"
)

func CreateSqlClient(connString *ConnectionString) (SqlUserClient, error) {

	var err error
	var db *sql.DB

	if connString.Username != "" {
		db, err = createUsingPasswordAuth(connString)
	} else {
		db, err = createUsingAzureActiveDirectoryAuth(connString)
	}

	if err != nil {
		return nil, err
	}

	dbId := parseDatabaseId(connString)

	return &sqlUserClient{conn: db, dbId: dbId}, nil
}

func createUsingPasswordAuth(connString *ConnectionString) (*sql.DB, error) {
	str, err := connString.String()
	if err != nil {
		return nil, err
	}

	return sql.Open("mssql", str)
}

func createUsingAzureActiveDirectoryAuth(connString *ConnectionString) (*sql.DB, error) {
	str, err := connString.String()
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] Using connection string ", str, " dupa")

	connector, err := mssql.NewAccessTokenConnector(str, func() (string, error) {
		token, err := cli.GetTokenFromCLI("https://database.windows.net/")
		if err != nil {
			return "", err
		}

		log.Println("[DEBUG] Access token %s", token.AccessToken)

		return token.AccessToken, nil
	})

	if err != nil {
		return nil, err
	}

	return sql.OpenDB(connector), nil
}

func parseDatabaseId(connString *ConnectionString) string {
	return connString.ServerAddress + "/" + connString.Database
}

type SqlUserClient interface {
	DatabaseId() string

	Get(name string) (*SqlUser, error)
	Create(name, password string, roles []string) error
	ChangePassword(name, password string) error
	ChangeRoles(name string, grant, revoke []string) error
	Delete(name string) error

	Close() error
}

type SqlUser struct {
	Name  string
	Roles []string
}

type sqlUserClient struct {
	conn *sql.DB
	dbId string
}

func (client *sqlUserClient) DatabaseId() string {
	return client.dbId
}

func (client *sqlUserClient) Get(name string) (*SqlUser, error) {
	rows, err := client.conn.Query(`
		SELECT u.name, r.name FROM sys.database_principals u
		LEFT JOIN sys.database_role_members m on u.principal_id = m.member_principal_id
		LEFT JOIN sys.database_principals r on r.principal_id = m.role_principal_id
		WHERE u.name = @name`,
		sql.Named("name", name))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var role sql.NullString

	if !rows.Next() {
		return nil, rows.Err()
	}

	if err := rows.Scan(&name, &role); err != nil {
		return nil, err
	}

	var roles []string

	if role.Valid {
		roles = append(roles, role.String)
	}

	for rows.Next() {
		rows.Scan(&name, &role)
		roles = append(roles, role.String)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	user := SqlUser{Name: name, Roles: roles}
	return &user, nil
}

func (client *sqlUserClient) Create(name, password string, roles []string) error {
	var cmd strings.Builder
	fmt.Fprintf(&cmd, "CREATE USER %s WITH PASSWORD = '%s'\n", name, password)
	for _, role := range roles {
		fmt.Fprintf(&cmd, "ALTER ROLE %s ADD MEMBER %s\n", role, name)
	}

	_, err := client.conn.Exec(cmd.String())
	return err
}

func (client *sqlUserClient) ChangePassword(name, password string) error {
	_, err := client.conn.Exec(`
		EXEC('ALTER USER ' + QUOTENAME(@user) + ' WITH PASSWORD = '''  + @password + ''''); 
		`,
		sql.Named("user", name),
		sql.Named("password", password))
	return err
}

func (client *sqlUserClient) ChangeRoles(name string, grant, revoke []string) error {
	var cmd strings.Builder
	for _, role := range grant {
		fmt.Fprintf(&cmd, "ALTER ROLE %s ADD MEMBER %s\n", role, name)
	}
	for _, role := range revoke {
		fmt.Fprintf(&cmd, "ALTER ROLE %s DROP MEMBER %s\n", role, name)
	}

	_, err := client.conn.Exec(cmd.String())
	return err
}

func (client *sqlUserClient) Delete(name string) error {
	_, err := client.conn.Exec(`
		EXEC('DROP USER ' + QUOTENAME(@user));
		`,
		sql.Named("user", name))
	return err
}

func (client *sqlUserClient) Close() error {
	return client.conn.Close()
}
