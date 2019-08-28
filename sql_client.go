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
	Get(name string) (*SqlUser, error)
	Create(name, password string) error
	ChangePassword(name, password string) error
	Delete(name string) error

	Close() error
}

type SqlUser struct {
	name string
}

type sqlUserClient struct {
	conn *sql.DB
}

func (client *sqlUserClient) Close() error {
	return client.conn.Close()
}

func (client *sqlUserClient) Get(name string) (*SqlUser, error) {
	err := client.conn.QueryRow(`
		SELECT TOP 1 name
		FROM sys.database_principals
		WHERE type NOT IN ('A', 'G', 'R', 'X')
			AND sid IS NOT NULL
			AND name = :name`, sql.Named("name", name)).Scan(&name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &SqlUser{name: name}, nil
}

func (client *sqlUserClient) Create(name, password string) error {
	_, err := client.conn.Exec(`
		DECLARE @Sql NVARCHAR(MAX) = 'CREATE USER ' + QUOTENAME(:user) + ' WITH PASSWORD = '''  + :password + ''''; 
		EXEC(@Sql)`,
		sql.Named("user", name),
		sql.Named("password", password))
	return err
}

func (client *sqlUserClient) ChangePassword(name, password string) error {
	_, err := client.conn.Exec(`
		DECLARE @Sql NVARCHAR(MAX) = 'ALTER USER ' + QUOTENAME(:user) + ' WITH PASSWORD = '''  + :password + ''''; 
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
