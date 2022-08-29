package sql

import (
	"database/sql"
	"fmt"
)

func CreateSqlUserClient(client *SqlClient) SqlUserClient {
	return &sqlUserClient{
		conn: client.Db,
		dbId: client.Id,
	}
}

type SqlUserClient interface {
	Get(name string) (*SqlUser, error)
	Create(user *SqlUser) error
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
		SELECT u.[name], u.[type] 
		FROM sys.database_principals u
		WHERE u.[name] = @name AND u.[type] in ('E', 'S', 'U')`,
		sql.Named("name", name))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, rows.Err()
	}

	var userType string

	if err := rows.Scan(&name, &userType); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	user := SqlUser{
		Name:          name,
		ExternalLogin: userType == "E"}
	return &user, nil
}

func (client *sqlUserClient) Create(user *SqlUser) error {
	var cmd string

	if user.ExternalLogin {
		if user.Password != "" {
			return fmt.Errorf("cannot create external sql user with password")
		}

		cmd = fmt.Sprintf("CREATE USER [%s] FROM EXTERNAL PROVIDER", user.Name)
	} else {
		cmd = fmt.Sprintf("CREATE USER [%s] WITH PASSWORD = '%s'\n", user.Name, user.Password)
	}

	_, err := client.conn.Exec(cmd)
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
	_, err := client.conn.Exec(fmt.Sprintf("DROP USER [%s]", name))
	return err
}
