# sqlc-gen-c: SQLite experimental branch

This branch is an experimental C99/sqlite3 generator for sqlc, based on
`f132b10b4a113cc4088c0569cfc7db04d2cf1dff`. It has been tested with sqlc
v1.31.1 and is intentionally limited to the query types needed by Procmac.
It is not a general C plugin or a published upstream release.

## Supported subset

- SQLite `:exec` and `:one` queries with contiguous, non-null TEXT parameters.
- TEXT result columns (nullable or non-null) and non-null INTEGER results.
- ASCII SQL source; parameter values may be non-ASCII NUL-terminated strings.
- Generated functions prepare, bind, step, and finalize statements, and return
  SQLite status codes. `:exec` returns a changed-row count. `:one` returns
  `SQLITE_NOTFOUND` when no row exists and copies text into owned allocations;
  call the generated row's `_free` function after a successful read.
- Generated functions use the caller's `sqlite3 *` connection and do not start,
  commit, or roll back transactions. The application owns transaction policy.

Other databases, SQL commands, parameter/result types, nullable INTEGER results,
embedded NUL result strings and non-ASCII SQL source are rejected. `:one`
returns the first row and does not assert uniqueness. Keep queries and schema
under version control and test generated functions with your own transactions.

## Build and configure

Build the process plugin with Go 1.23 or newer:

```sh
go test ./...
go build -o sqlc-gen-c ./plugin
```

Point sqlc's version 2 configuration at the built plugin, for example:

```yaml
version: "2"
plugins:
  - name: c-sqlite
    process:
      cmd: ./sqlc-gen-c
sql:
  - engine: sqlite
    schema: schema.sql
    queries: queries.sql
    codegen:
      - plugin: c-sqlite
        out: generated
```

Then run `sqlc generate`. The C code links SQLite directly and does not require a
Go runtime. The original upstream WASM-oriented end-to-end fixtures are kept in
`testdata/upstream-legacy-endtoend/`; they do not test this process-plugin fork.
See `internal/sqlite_test.go` for the fork's generator tests. MIT licensed; see
`LICENSE`.
