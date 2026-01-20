package c

import (
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

type Field struct {
	Name    string
	DBName  string // Name as used in the DB
	Type    string
	Comment string
	Column  *plugin.Column
	// EmbedFields contains the embedded fields that require scanning.
	EmbedFields []Field
}

func (gf Field) HasSqlcSlice() bool {
	return gf.Column.IsSqlcSlice
}
