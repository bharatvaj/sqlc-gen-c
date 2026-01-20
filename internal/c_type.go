package c

import (
	"strings"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/timwmillard/sqlc-gen-c/internal/opts"
)

func cType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	typ := cInnerType(req, options, col)
	if col.IsSqlcSlice {
		return typ + "[]"
	}
	if col.IsArray {
		return typ + strings.Repeat("[]", int(col.ArrayDims))
	}
	return typ
}

func cInnerType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	switch req.Settings.Engine {
	case "mysql":
		return mysqlType(req, options, col)
	case "postgresql":
		return postgresType(req, options, col)
	case "sqlite":
		return sqliteType(req, options, col)
	default:
		return "void*"
	}
}
