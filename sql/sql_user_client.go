package sql

import (
	"database/sql"
	"fmt"
	"strings"
)

func CreateSqlUserClient(client *SqlClient) SqlUserClient {
	return &sqlUserClient{
		conn: client.Db,
		dbId: client.Id,
	}
}

func parseDatabaseId(connString *ConnectionString) string {
	return connString.ServerAddress + "/" + connString.Database
}

type SqlUserClient interface {
	Get(name string) (*SqlUser, error)
	Create(name, password string) error
	ChangePassword(name, password string) error
	Delete(name string) error
}

type SqlUser struct {
	Name          string
	Password      string
	ExternalLogin bool
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
		SELECT u.name FROM sys.database_principals u
		WHERE u.name = @name`,
		sql.Named("name", name))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, rows.Err()
	}

	if err := rows.Scan(&name); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	user := SqlUser{Name: name}
	return &user, nil
}

func (client *sqlUserClient) Create(name, password string) error {
	var cmd strings.Builder
	fmt.Fprintf(&cmd, "CREATE USER %s WITH PASSWORD = '%s'\n", name, password)

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

func (client *sqlUserClient) Delete(name string) error {
	_, err := client.conn.Exec(`
		EXEC('DROP USER ' + QUOTENAME(@user));
		`,
		sql.Named("user", name))
	return err
}
