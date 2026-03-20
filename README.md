# GO-RSS-AGG
Project in Go for RSS Aggregation


Initialize module using `go mod init`.
Install packages using `go get <package-name>`.

Added goose and sqlc for postgres migrations and conversion of sql queries into type-safe Go code respectively.

Define `sql/` directory for holding queries, schemas etc.
Define sqlc configuration with `sqlc.yaml`.

Run sqlc commands from project root. Use the `generate` command to generate Go code from SQL queries.

`sqlc.yaml` defines where to output the generated Go code: `internal/database/`.

Import this internal package into main.

Set up database connection using `database/sql` with the `lib/pq` postgres driver, then pass it to the sqlc-generated `Queries` struct via a `dbConfig` wrapper.

Added a user creation endpoint (`POST /v1/users`) with a model conversion layer (`models.go`) to map between sqlc-generated types and our JSON response types.

