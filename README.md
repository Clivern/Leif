## Leif

Stdio MCP server for local PostgreSQL. Tools: `list_tables`, `describe_table`, `query`.

```bash
go run . serve --dsn 'postgres://user@127.0.0.1:5432/dbname?sslmode=disable'
```

Or set `DATABASE_URL` and run `go run . serve`.

### Per-project Cursor config

`.cursor/mcp.json` in the app repo:

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

On `PATH`: `"command": "leif"`, `"args": ["serve"]`. Reload, then enable it in **Settings → MCP**.

Workspace `.cursor/mcp.json` overrides `~/.cursor/mcp.json`. To keep the DSN out of git, use `"DATABASE_URL": "${env:DATABASE_URL}"`.

### Ask Cursor

> Use Leif to list tables, then describe `public.orders`.

Optional `.cursor/rules/leif.mdc` so it does this without being asked:

```markdown
---
description: Use Leif for PostgreSQL questions
alwaysApply: true
---

Use Leif MCP tools for schema and data. Call `list_tables` or `describe_table` before SQL, then `query`.
```
