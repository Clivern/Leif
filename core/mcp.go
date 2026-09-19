// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server represents the PostgreSQL MCP server.
type Server struct {
	db *sql.DB
}

// QueryArgs represents the arguments for the query tool.
type QueryArgs struct {
	SQL string `json:"sql" jsonschema:"SQL to run"`
}

// DescribeArgs represents the arguments for the describe tool.
type DescribeArgs struct {
	Schema string `json:"schema,omitempty" jsonschema:"schema name, defaults to public"`
	Table  string `json:"table" jsonschema:"table name"`
}

// Serve starts the PostgreSQL MCP server over stdio.
func Serve(ctx context.Context, db *sql.DB, version string) error {
	s := &Server{db: db}

	mcpServer := mcp.NewServer(&mcp.Implementation{Name: "leif", Version: version}, nil)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_tables",
		Description: "List tables in the connected PostgreSQL database",
	}, s.ListTables)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "describe_table",
		Description: "Describe columns for a PostgreSQL table",
	}, s.DescribeTable)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "query",
		Description: "Run SQL against the connected PostgreSQL database",
	}, s.Query)

	return mcpServer.Run(ctx, &mcp.StdioTransport{})
}

// ListTables lists the tables in the connected PostgreSQL database.
func (s *Server) ListTables(context.Context, *mcp.CallToolRequest, struct{}) (*mcp.CallToolResult, any, error) {
	rows, err := s.db.Query(`
		SELECT table_schema, table_name
		FROM information_schema.tables
		WHERE table_type = 'BASE TABLE'
			AND table_schema NOT IN ('pg_catalog', 'information_schema')
		ORDER BY table_schema, table_name`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	return JsonResult(ScanMaps(rows))
}

// DescribeTable describes the columns for a PostgreSQL table.
func (s *Server) DescribeTable(_ context.Context, _ *mcp.CallToolRequest, in DescribeArgs) (*mcp.CallToolResult, any, error) {
	schema := in.Schema
	if schema == "" {
		schema = "public"
	}

	rows, err := s.db.Query(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2
		ORDER BY ordinal_position`, schema, in.Table)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	return JsonResult(ScanMaps(rows))
}

// Query runs SQL against the connected PostgreSQL database.
func (s *Server) Query(_ context.Context, _ *mcp.CallToolRequest, in QueryArgs) (*mcp.CallToolResult, any, error) {
	rows, err := s.db.Query(in.SQL)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	return JsonResult(ScanMaps(rows))
}

// ScanMaps scans the rows into a slice of maps.
func ScanMaps(rows *sql.Rows) []map[string]any {
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	var out []map[string]any
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}

		err = rows.Scan(ptrs...)
		if err != nil {
			panic(err)
		}

		row := map[string]any{}
		for i, col := range cols {
			val := vals[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			row[col] = val
		}
		out = append(out, row)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return out
}

// JsonResult returns the result as a JSON string.
func JsonResult(v any) (*mcp.CallToolResult, any, error) {
	payload, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(payload)}},
	}, nil, nil
}
