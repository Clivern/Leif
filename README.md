## Leif

Stdio MCP server for a local PostgreSQL database.

```bash
go run . serve --dsn 'postgres://user@127.0.0.1:5432/dbname?sslmode=disable'
```

Or set `DATABASE_URL` and run `go run . serve`.

Leif exposes three tools: `list_tables`, `describe_table`, and `query`.

### Configure it per project

Add a project-level MCP config so only that repo's database is available.

1. In the app repo (not necessarily this one), create `.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "leif": {
      "command": "go",
      "args": ["run", "/absolute/path/to/leif", "serve"],
      "env": {
        "DATABASE_URL": "postgres://user@127.0.0.1:5432/dbname?sslmode=disable"
      }
    }
  }
}
```

If Leif is already on your `PATH`, use `"command": "leif"` and `"args": ["serve"]` instead of `go run`.

2. Point `DATABASE_URL` at **this** project's Postgres. Prefer `${env:DATABASE_URL}` if you do not want the DSN in git:

```json
"env": {
  "DATABASE_URL": "${env:DATABASE_URL}"
}
```

3. Reload the window (or restart Cursor), then open **Settings → MCP** and confirm `leif` is enabled and connected. The first time Cursor starts it, approve the server if prompted.

Project config (`.cursor/mcp.json`) applies only to that workspace. A global `~/.cursor/mcp.json` is available in every project; if both define `leif`, the project file wins.

### Ask Cursor to use it

Agent chat can call Leif once the server is connected. Be explicit about the database:

- "Use Leif to list the tables in this database."
- "Ask Leif to describe `public.orders`, then write a query for unpaid rows."
- "Query the database with Leif — do not guess the schema."

Cursor may ask you to approve each tool call (`list_tables`, `describe_table`, `query`). Allow those when the SQL looks right.

To make this automatic, add a project rule at `.cursor/rules/leif.mdc`:

```markdown
---
description: Use Leif for PostgreSQL questions
alwaysApply: true
---

For schema or data questions, use the Leif MCP tools instead of guessing.
Call `list_tables` or `describe_table` before writing SQL, then `query`.
```
