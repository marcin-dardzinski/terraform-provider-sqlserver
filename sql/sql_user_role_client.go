package sql

import (
	"database/sql"
	"fmt"
)

type SqlUserRole struct {
	User string
	Role string
}

type SqlUserRoleClient interface {
	Exists(role SqlUserRole) (bool, error)
	Grant(role SqlUserRole) error
	Revoke(role SqlUserRole) error
}

type sqlUserRoleClient struct {
	conn *sql.DB
	dbId string
}

func CreateSqlUserRoleClient(client *SqlClient) SqlUserRoleClient {
	return &sqlUserRoleClient{
		conn: client.Db,
		dbId: client.Id,
	}
}

func (client *sqlUserRoleClient) Exists(role SqlUserRole) (bool, error) {
	row := client.conn.QueryRow(`
		select top 1 1
		from sys.database_principals r
		join sys.database_role_members m on m.role_principal_id = r.principal_id
		join sys.database_principals u on m.member_principal_id = u.principal_id
		where r.[name] = @rolename and r.[type] = 'R' and u.[name] = @username and u.[type] in ('C', 'E', 'K', 'S', 'U') `,
		sql.Named("username", role.User),
		sql.Named("rolename", role.Role))

	var placeholder int

	err := row.Scan(&placeholder)

	if err == nil {
		return true, nil
	} else if err == sql.ErrNoRows {
		return false, nil
	} else {
		return false, err
	}

}

func (client *sqlUserRoleClient) Grant(role SqlUserRole) error {
	_, err := client.conn.Exec(
		fmt.Sprintf("ALTER ROLE %s ADD MEMBER %s", role.Role, role.User))

	return err
}

func (client *sqlUserRoleClient) Revoke(role SqlUserRole) error {
	_, err := client.conn.Exec(
		fmt.Sprintf("ALTER ROLE %s DROP MEMBER %s", role.Role, role.User))

	return err
}
