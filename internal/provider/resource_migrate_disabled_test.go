package provider

import (
	"fmt"
	"testing"

	helperresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceNullMigrate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test")
	}

	for _, server := range testServers {
		t.Run(server.ServerType, func(t *testing.T) {
			url, _, err := server.URL()
			if err != nil {
				t.Fatal(err)
			}

			helperresource.UnitTest(t, helperresource.TestCase{
				ProtoV6ProviderFactories: protoV6ProviderFactories,
				Steps: []helperresource.TestStep{
					{
						Config: fmt.Sprintf(`
provider "sql" {
	url 	 = null # %q
	disabled = true

	max_idle_conns = 0
}

resource "sql_migrate" "db" {
	count = 0

	migration {
		id = "create table"

		up = <<SQL
CREATE TABLE inline_migrate_test (
	user_id integer unique,
	name    varchar(40),
	email   varchar(40)
);
SQL

		down = <<SQL
DROP TABLE inline_migrate_test;
SQL
	}
}

data "sql_query" "users" {
	count 	   = 0
	depends_on = [sql_migrate.db]

	query = "select * from inline_migrate_test"
}

output "rowcount" {
	value = length(data.sql_query.users)
}
				`, url),
						Check: helperresource.ComposeTestCheckFunc(
							helperresource.TestCheckOutput("rowcount", "0"),
						),
					},
					{
						Config: fmt.Sprintf(`
provider "sql" {
	url      = null # %q
	disabled = true

	max_idle_conns = 0
}

resource "sql_migrate" "db" {
	count = 0

	migration {
		id = "create table"

		up = <<SQL
CREATE TABLE inline_migrate_test (
	user_id integer unique,
	name    varchar(40),
	email   varchar(40)
);
SQL

		down = <<SQL
DROP TABLE inline_migrate_test;
SQL
	}

	migration {
		id   = "insert row"
		up   = "INSERT INTO inline_migrate_test VALUES (1, 'Paul Tyng', 'paul@example.com');"
		down = "DELETE FROM inline_migrate_test WHERE user_id = 1;"
	}
}

data "sql_query" "users" {
	count 	   = 0
	depends_on = [sql_migrate.db]

	query = "select * from inline_migrate_test"
}

output "rowcount" {
	value = length(data.sql_query.users)
}
				`, url),
						Check: helperresource.ComposeTestCheckFunc(
							helperresource.TestCheckOutput("rowcount", "0"),
						),
					},
				},
			})
		})
	}
}
