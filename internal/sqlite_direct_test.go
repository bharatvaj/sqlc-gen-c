package c

import (
	"context"
	"strings"
	"testing"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

func column(name, typ string, notNull bool) *plugin.Column {
	return &plugin.Column{Name: name, Type: &plugin.Identifier{Name: typ}, NotNull: notNull}
}

func TestSQLiteGenerator(t *testing.T) {
	good := &plugin.Query{Name: "Get", Cmd: ":one", Text: "SELECT value FROM item WHERE id = ?1",
		Params:  []*plugin.Parameter{{Number: 1, Column: column("id", "TEXT", true)}},
		Columns: []*plugin.Column{column("value", "TEXT", false)}}
	req := &plugin.GenerateRequest{Settings: &plugin.Settings{Engine: "sqlite"}, PluginOptions: []byte(`{"sqlite_query_api":"direct"}`), Queries: []*plugin.Query{good}}
	out, err := Generate(context.Background(), req)
	if err != nil || len(out.GetFiles()) != 2 {
		t.Fatalf("generate: %v", err)
	}
	code := string(out.Files[1].Contents)
	for _, required := range []string{"sqlite3_bind_text(stmt, 1, p1", "SQLITE_NOTFOUND", "sqlite3_finalize(stmt)", "copy_text(stmt, 0, 1", "sqlc_Get_free(&row)"} {
		if !strings.Contains(code, required) {
			t.Errorf("generated code lacks %q", required)
		}
	}
}

func TestSQLiteRejectsUnsupported(t *testing.T) {
	cases := []*plugin.Query{
		{Name: "Many", Cmd: ":many", Text: "SELECT 1"},
		{Name: "Blob", Cmd: ":one", Text: "SELECT x", Columns: []*plugin.Column{column("x", "BLOB", true)}},
		{Name: "NullableInteger", Cmd: ":one", Text: "SELECT x", Columns: []*plugin.Column{column("x", "INTEGER", false)}},
		{Name: "NumberParam", Cmd: ":exec", Text: "SELECT ?1", Params: []*plugin.Parameter{{Number: 1, Column: column("p", "INTEGER", true)}}},
		{Name: "Gap", Cmd: ":exec", Text: "SELECT ?2", Params: []*plugin.Parameter{{Number: 2, Column: column("p", "TEXT", true)}}},
		{Name: "Unsafe", Cmd: ":one", Text: "SELECT x", Columns: []*plugin.Column{column("int", "TEXT", true)}},
	}
	for _, q := range cases {
		_, err := Generate(context.Background(), &plugin.GenerateRequest{Settings: &plugin.Settings{Engine: "sqlite"}, PluginOptions: []byte(`{"sqlite_query_api":"direct"}`), Queries: []*plugin.Query{q}})
		if err == nil {
			t.Errorf("%s: expected rejection", q.Name)
		}
	}
	if _, err := Generate(context.Background(), &plugin.GenerateRequest{Settings: &plugin.Settings{Engine: "postgresql"}, PluginOptions: []byte(`{"sqlite_query_api":"direct"}`)}); err == nil {
		t.Error("expected non-SQLite rejection")
	}
}
