provider "sqlserver" {
  connection_string = "Server=localhost;Database=Db1;User Id=sa;Password=Passwd1!;"
}

resource "sqlserver_user" "main" {
  name     = "foo55"
  password = "Passwd1!2"
  roles = [
    "db_datareader",
    "db_datawriter"
  ]
}

# Server=localhost;Database=Db1;User Id=sa;Password=Passwd1!;
