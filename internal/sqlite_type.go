package c

import (
	"log"
	"strings"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
	"github.com/timwmillard/sqlc-gen-c/internal/debug"
	"github.com/timwmillard/sqlc-gen-c/internal/opts"
)

func sqliteType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	dt := strings.ToLower(sdk.DataType(col.Type))
	notNull := col.NotNull || col.IsArray
	emitPointersForNull := options.EmitPointersForNullTypes

	switch dt {

	case "int", "integer", "tinyint", "smallint", "mediumint", "bigint", "unsignedbigint", "int2", "int8":
		if notNull {
			return "sql_int64"
		}
		if emitPointersForNull {
			return "sql_int64*"
		}
		return "sql_nullint64"

	case "blob":
		return "[]byte"

	case "real", "double", "doubleprecision", "float":
		if notNull {
			return "sql_double"
		}
		if emitPointersForNull {
			return "sql_double*"
		}
		return "sql_nulldouble"

	case "boolean", "bool":
		if notNull {
			return "sql_bool"
		}
		if emitPointersForNull {
			return "sql_bool*"
		}
		return "sql_nullbool"

	case "date", "datetime", "timestamp":
		if notNull {
			return "time.Time"
		}
		if emitPointersForNull {
			return "*time.Time"
		}
		return "sql.NullTime"

	case "any":
		return "void*"

	}

	switch {

	case strings.HasPrefix(dt, "character"),
		strings.HasPrefix(dt, "varchar"),
		strings.HasPrefix(dt, "varyingcharacter"),
		strings.HasPrefix(dt, "nchar"),
		strings.HasPrefix(dt, "nativecharacter"),
		strings.HasPrefix(dt, "nvarchar"),
		dt == "text",
		dt == "clob":
		if notNull {
			return "sql_text"
		}
		if emitPointersForNull {
			return "sql_text*"
		}
		return "sql_nulltext"

	case strings.HasPrefix(dt, "decimal"), dt == "numeric":
		if notNull {
			return "sql_numeric"
		}
		if emitPointersForNull {
			return "sql_numeric*"
		}
		return "sql_nullnumeric"

	default:
		if debug.Active {
			log.Printf("unknown SQLite type: %s\n", dt)
		}

		return "void*"

	}
}
