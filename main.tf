provider "sqlserver" {
  # connection_string = "Server=localhost;Port=1433;Database=Db1;User Id=sa;Password=Passwd1!;"
  server   = "localhost"
  database = "Db1"
  username = "sa"
  password = "Passwd1!"
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

data "sqlserver_connection_string" "root" {
  server   = "localhost"
  database = "foo"
  username = "foo"
  password = "foo"
}

output "foo" {
  value     = data.sqlserver_connection_string.root.value
  sensitive = true
}

