terraform {
  required_providers {
    sqlserver = {
      source  = "marcin-dardzinski/sqlserver"
      version = "0.0.1"
    }
  }
}

provider "sqlserver" {
  connection_string = "Server=localhost,1433;Database=tf-tests;Persist Security Info=False;User ID=sa;Password=Passwd1!;"
}

resource "sqlserver_user" "foo55" {
  name     = "foo55"
  password = "Passwd1!2"
}

resource "sqlserver_user_role" "foo" {
  user = sqlserver_user.foo55.name
  role = "db_ddladmin"
}
