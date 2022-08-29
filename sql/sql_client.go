package sql

import (
	"database/sql"

	"github.com/denisenkom/go-mssqldb/azuread"
	"github.com/denisenkom/go-mssqldb/msdsn"
)

type SqlClientConfig struct {
	ConnectionString string
}

type SqlClient struct {
	Db *sql.DB
	Id string
}

func CreateSqlClient(config SqlClientConfig) (*SqlClient, error) {
	parsedConnStr, _, err := msdsn.Parse(config.ConnectionString)

	if err != nil {
		return nil, err
	}

	db, err := sql.Open(azuread.DriverName, config.ConnectionString)

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
