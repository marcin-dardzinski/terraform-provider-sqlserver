provider "sqlserver" {
  connection_string = "Server=localhost;Database=Db1;User Id=sa;Password=Passwd1!;"
  database_id       = "my-awesome-db-id"
}

resource "sqlserver_user" "foo55" {
  name     = "foo55"
  password = "Passwd1!2"
  roles = [
    "db_datareader",
    "db_datawriter",
    "db_ddladmin"
  ]
}

resource "sqlserver_user" "foo" {
  name     = "foo"
  password = "Passwd1!2"
  roles = [
    "db_datareader",
    "db_datawriter"
  ]
}
