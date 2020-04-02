# Templating for Migration Sources

migrate-template a wrapper for migration sources of library [golang-migrate/migrate](https://github.com/golang-migrate/migrate) that allows the use of variables in migrations files.



## Use in your Go project

```go

import (
	template "github.com/FindMyKids/migrate-template"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
)

func main() {
	sourceInstance, err := source.Open("file://migration")
	
	templateInstance := template.Wrap(
		sourceInstance,
		migration.WithVars(template.M{
			"cluster": "cluster_name",
			"replicated_path": "/clickhouse/tables/db_name",
		}),
	)

	m, err := migrate.NewWithSourceInstance(
		"template", templateInstance,
		"clickhouse://...",
	)
	m.Steps(1)
}

```

Migration template using variables:
```sql

CREATE TABLE table_name ON CLUSTER {{cluster}} (
	foo Int64,
	bar String
) ENGINE = ReplicatedReplacingMergeTree('{{replicated_path}}.table_name', '{replica}')
```

Result:
```sql

CREATE TABLE table_name ON CLUSTER cluster_name (
	foo Int64,
	bar String
) ENGINE = ReplicatedReplacingMergeTree('/clickhouse/tables/db_name.table_name', '{replica}')
```