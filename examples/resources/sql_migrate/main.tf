terraform {
  required_providers {
    sql = {
      source = "divisionium/sql"
    }
  }
}

provider "sql" {
  url = "postgres://postgres:tf@localhost:5432/tftest?sslmode=disable"
}
