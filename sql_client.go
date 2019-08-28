package main

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

func CreateSqlClient(connString string) (SqlUserClient, error) {
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, err
	}

	return &sqlUserClient{conn: conn}, nil
}

type SqlUserClient interface {
	Exists(name string) (bool, error)
	Create(name, password string) error
	Delete(name string) error

	Close() error
}

type sqlUserClient struct {
	conn *sql.DB
}

func (client *sqlUserClient) Close() error {
	return client.conn.Close()
}

func (client *sqlUserClient) Exists(name string) (bool, error) {
	row, err := client.conn.Query(`
		SELECT TOP 1 1
		FROM sys.database_principals
		WHERE type NOT IN ('A', 'G', 'R', 'X')
			AND sid IS NOT NULL
			AND name = :name`, sql.Named("name", name))

	if err != nil {
		return false, err
	}

	defer row.Close()

	return row.Next(), nil
}

func (client *sqlUserClient) Create(name, password string) error {
	_, err := client.conn.Exec(`
		DECLARE @Sql NVARCHAR(MAX) = 'CREATE USER ' + QUOTENAME(:user) + ' WITH PASSWORD = '''  + :password + ''''; 
		EXEC(@Sql)`,
		sql.Named("user", name),
		sql.Named("password", password))
	return err
}

func (client *sqlUserClient) Delete(name string) error {
	_, err := client.conn.Exec(`
		DECLARE @Sql NVARCHAR(MAX) = 'DROP USER ' + QUOTENAME(:user);
		EXEC(@Sql)`,
		sql.Named("user", name))
	return err
}
