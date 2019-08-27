package main

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

type SqlClientProvider interface {
	GetClient() (SqlUserClient, error)
}

type SqlUserClient interface {
	Exists(name string) (bool, error)
	Create(name, password string) error
	Delete(name string) error

	Close() error
}

type sqlClientProvider struct {
	connString string
}

type sqlUserClient struct {
	conn *sql.DB
}

func (provider *sqlClientProvider) GetClient() (SqlUserClient, error) {
	conn, err := sql.Open("mssql", provider.connString)
	if err != nil {
		return nil, err
	}

	client := sqlUserClient{
		conn: conn,
	}

	return &client, nil
}

func (client *sqlUserClient) Close() error {
	return client.conn.Close()
}

func (client *sqlUserClient) Exists(name string) (bool, error) {
	// return false, nil
	row, err := client.conn.Query(`
		SELECT TOP 1 1
		FROM sys.database_principals
		WHERE type NOT IN ('A', 'G', 'R', 'X')
			AND sid IS NOT NULL
			AND name = @name`, sql.Named("name", name))

	if err != nil {
		return false, err
	}

	defer row.Close()

	return row.Next(), nil
}

func (client *sqlUserClient) Create(name, password string) error {
	// _, err := client.conn.Exec(
	// 	`EXEC( 'CREATE USER ' + QUOTENAME(@user) + ' WITH PASSWORD = ' @password )`,
	// 	sql.Named("user", name),
	// 	sql.Named("password", password))
	// return err
	_, err := client.conn.Exec(`CREATE USER foo2 WITH PASSWORD = 'Passwd1!'`)
	return err
}

func (client *sqlUserClient) Delete(name string) error {
	_, err := client.conn.Exec(`EXEC( 'DELETE USER ' + QUOTENAME(@user))`, sql.Named("user", name))
	return err
}
