package opts

import (
	"encoding/json"
	"fmt"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

// version: 2
// plugins:
// - name: c
//   wasm:
//     url: file:///Users/tim/go/src/sqlc-gen-c/bin/sqlc-gen-c.wasm
//     sha256: ""
// sql:
// - engine: "sqlite"
//   queries: "queries.sql"
//   schema: "schema.sql"
//   codegen:
//   - plugin: c
//     out: gen
//     options:
//       prefix: sql_
//       sql_type_case: snake
//       struct_type_case: snake
//       field_type_case: snake
//       func_type_case: snake
//       sys_include:
//         stdio.h
//       include:
//         sqlite.h
//       output_h: models.h
//       output_c: models.c
//       sql_context:
//         - void *(*alloc)(void *, int)
//         - areana a
//       pointers_for_null_types: false
// 		 emit_sql_as_comment: true

type Options struct {
	// New
	SQLiteQueryAPI string   `json:"sqlite_query_api,omitempty" yaml:"sqlite_query_api"`
	Prefix         string   `json:"prefix,omitempty" yaml:"prefix"`
	SQLTypeCase    string   `json:"sql_type_case,omitempty" yaml:"sql_type_case"`
	StructTypeCase string   `json:"struct_type_case,omitempty" yaml:"struct_type_case"`
	FieldTypeCase  string   `json:"field_type_case,omitempty" yaml:"field_type_case"`
	FuncTypeCase   string   `json:"func_type_case,omitempty" yaml:"func_type_case"`
	SysInclude     []string `json:"sys_include,omitempty" yaml:"sys_include"`
	Include        []string `json:"include,omitempty" yaml:"include"`
	OutputH        string   `json:"output_h,omitempty" yaml:"output_h"`
	OutputC        string   `json:"output_c,omitempty" yaml:"output_c"`
	SQLContext     []string `json:"sql_context,omitempty" yaml:"sql_context"`

	// Existing
	EmitPointersForNullTypes bool      `json:"emit_pointers_for_null_types,omitempty" yaml:"emit_pointers_for_null_types"`
	EmitSQLAsComment         bool      `json:"emit_sql_as_comment,omitempty" yaml:"emit_sql_as_comment"`
	EmitExportedQueries      bool      `json:"emit_exported_queries" yaml:"emit_exported_queries"`
	Out                      string    `json:"out,omitempty" yaml:"out"`
	SQLDriver                string    `json:"sql_driver,omitempty" yaml:"sql_driver"`
	OmitUnusedStructs        bool      `json:"omit_unused_structs,omitempty" yaml:"omit_unused_structs"`
	Initialisms              *[]string `json:"initialisms,omitempty" yaml:"initialisms"`
	EmitResultStructPointers bool      `json:"emit_result_struct_pointers" yaml:"emit_result_struct_pointers"`
	EmitParamsStructPointers bool      `json:"emit_params_struct_pointers" yaml:"emit_params_struct_pointers"`

	InitialismsMap map[string]struct{} `json:"-,omitempty" yaml:"-"`

	// MAYBE
	EmitExactTableNames bool   `json:"emit_exact_table_names,omitempty" yaml:"emit_exact_table_names"`
	EmitEmptySlices     bool   `json:"emit_empty_slices,omitempty" yaml:"emit_empty_slices"`
	QueryParameterLimit *int32 `json:"query_parameter_limit,omitempty" yaml:"query_parameter_limit"`
	OmitSqlcVersion     bool   `json:"omit_sqlc_version,omitempty" yaml:"omit_sqlc_version"`
}

func Parse(req *plugin.GenerateRequest) (*Options, error) {
	options, err := parseOpts(req)
	if err != nil {
		return nil, err
	}
	return options, nil
}

func parseOpts(req *plugin.GenerateRequest) (*Options, error) {
	var options Options
	if len(req.PluginOptions) == 0 {
		return &options, nil
	}
	if err := json.Unmarshal(req.PluginOptions, &options); err != nil {
		return nil, fmt.Errorf("unmarshalling plugin options: %w", err)
	}

	if options.QueryParameterLimit == nil {
		options.QueryParameterLimit = new(int32)
		*options.QueryParameterLimit = 1
	}

	if options.Initialisms == nil {
		options.Initialisms = new([]string)
		*options.Initialisms = []string{"id"}
	}

	options.InitialismsMap = map[string]struct{}{}
	for _, initial := range *options.Initialisms {
		options.InitialismsMap[initial] = struct{}{}
	}

	return &options, nil
}

func ValidateOpts(opts *Options) error {
	if *opts.QueryParameterLimit < 0 {
		return fmt.Errorf("invalid options: query parameter limit must not be negative")
	}

	return nil
}
