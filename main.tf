provider "sqlserver" {
  connection_string = "Server=localhost;Database=Db1;User Id=sa;Password=Passwd1!;"
}

resource "sqlserver_user" "main" {
  name     = "foo5"
  password = "Passwd1!"
}

# Server=localhost;Database=Db1;User Id=sa;Password=Passwd1!;
