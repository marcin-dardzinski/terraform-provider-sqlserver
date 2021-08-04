terraform {
  required_providers {
    sqlserver = {
      source  = "local/local/sqlserver"
      version = "1.0.0"
    }
  }
}

provider "sqlserver" {
  connection_string = "Server=tf-2137.database.windows.net;Port=1433;Database=tf-2137;"
  azure {
  }
}

resource "sqlserver_user" "foo55" {
  name     = "foo55"
  password = "Passwd1!2"
}

resource "sqlserver_user" "foo" {
  name     = "foo"
  password = "Passwd1!2"
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

