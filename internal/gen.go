package c

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
	"github.com/timwmillard/sqlc-gen-c/internal/opts"
)

type tmplCtx struct {
	Enums       []Enum
	Structs     []Struct
	Queries     []Query
	SqlcVersion string

	// TODO: Race conditions
	SourceName string
	Options    opts.Options
}

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	options, err := opts.Parse(req)
	if err != nil {
		return nil, err
	}

	if err := opts.ValidateOpts(options); err != nil {
		return nil, err
	}

	enums := buildEnums(req, options)
	structs := buildStructs(req, options)
	queries, err := buildQueries(req, options, structs)
	if err != nil {
		return nil, err
	}

	if options.OmitUnusedStructs {
		enums, structs = filterUnusedStructs(enums, structs, queries)
	}

	if err := validate(options, enums, structs, queries); err != nil {
		return nil, err
	}

	return generate(req, options, enums, structs, queries)
}

func validate(options *opts.Options, enums []Enum, structs []Struct, queries []Query) error {
	enumNames := make(map[string]struct{})
	for _, enum := range enums {
		enumNames[enum.Name] = struct{}{}
		enumNames["Null"+enum.Name] = struct{}{} // TODO: Null casing
	}
	structNames := make(map[string]struct{})
	for _, struckt := range structs {
		if _, ok := enumNames[struckt.Name]; ok {
			return fmt.Errorf("struct name conflicts with enum name: %s", struckt.Name)
		}
		structNames[struckt.Name] = struct{}{}
	}
	for _, query := range queries {
		if _, ok := enumNames[query.ConstantName]; ok {
			return fmt.Errorf("query constant name conflicts with enum name: %s", query.ConstantName)
		}
		if _, ok := structNames[query.ConstantName]; ok {
			return fmt.Errorf("query constant name conflicts with struct name: %s", query.ConstantName)
		}
	}
	return nil
}

func generate(req *plugin.GenerateRequest, options *opts.Options, enums []Enum, structs []Struct, queries []Query) (*plugin.GenerateResponse, error) {
	tctx := tmplCtx{
		Enums:       enums,
		Structs:     structs,
		SqlcVersion: req.SqlcVersion,
		Queries:     queries,
		Options:     *options,
	}

	funcMap := template.FuncMap{
		"lowerTitle": sdk.LowerTitle,
		"comment":    sdk.DoubleSlashComment,
		"escape":     escapeCString,
	}

	tmpl := template.Must(
		template.New("table").
			Funcs(funcMap).
			ParseFS(
				templates,
				"templates/*.tmpl",
				"templates/*/*.tmpl",
			),
	)

	output := map[string]string{}

	execute := func(name, templateName string) error {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		err := tmpl.ExecuteTemplate(w, templateName, &tctx)
		w.Flush()
		if err != nil {
			return err
		}

		if !strings.HasSuffix(name, ".c") && !strings.HasSuffix(name, ".h") {
			name += ".c"
		}
		output[name] = b.String()
		return nil
	}

	if req.Settings.Engine == "sqlite" {
		if err := execute("models.h", "sqliteModel"); err != nil {
			return nil, err
		}
	}
	if req.Settings.Engine == "postgres" {
		if err := execute("model_psql.c", "libpqModel"); err != nil {
			return nil, err
		}
	}

	files := map[string]struct{}{}
	for _, gq := range queries {
		files[gq.SourceName] = struct{}{}
	}

	for source := range files {
		if err := execute(source, "queryFile"); err != nil {
			return nil, err
		}
	}
	resp := plugin.GenerateResponse{}

	for filename, code := range output {
		resp.Files = append(resp.Files, &plugin.File{
			Name:     filename,
			Contents: []byte(code),
		})
	}

	// debugOutput := fmt.Sprintf("queries=%d\n", len(queries))
	// resp.Files = append(resp.Files, &plugin.File{
	// 	Name:     "output.txt",
	// 	Contents: []byte(debugOutput),
	// })

	return &resp, nil
}

func filterUnusedStructs(enums []Enum, structs []Struct, queries []Query) ([]Enum, []Struct) {
	keepTypes := make(map[string]struct{})

	for _, query := range queries {
		if !query.Arg.isEmpty() {
			keepTypes[query.Arg.Type()] = struct{}{}
			if query.Arg.IsStruct() {
				for _, field := range query.Arg.Struct.Fields {
					keepTypes[field.Type] = struct{}{}
				}
			}
		}
		if query.hasRetType() {
			keepTypes[query.Ret.Type()] = struct{}{}
			if query.Ret.IsStruct() {
				for _, field := range query.Ret.Struct.Fields {
					keepTypes[field.Type] = struct{}{}
					for _, embedField := range field.EmbedFields {
						keepTypes[embedField.Type] = struct{}{}
					}
				}
			}
		}
	}

	keepEnums := make([]Enum, 0, len(enums))
	for _, enum := range enums {
		_, keep := keepTypes[enum.Name]
		_, keepNull := keepTypes["Null"+enum.Name]
		if keep || keepNull {
			keepEnums = append(keepEnums, enum)
		}
	}

	keepStructs := make([]Struct, 0, len(structs))
	for _, st := range structs {
		if _, ok := keepTypes[st.Name]; ok {
			keepStructs = append(keepStructs, st)
		}
	}

	return keepEnums, keepStructs
}

func escapeCString(s string) string {
	// Remove trailing newline to avoid empty last line
	s = strings.TrimSuffix(s, "\n")
	lines := strings.Split(s, "\n")
	var result strings.Builder

	for i, line := range lines {
		// Escape special characters
		escaped := strings.ReplaceAll(line, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")

		result.WriteString("                      \"")
		result.WriteString(escaped)
		result.WriteString("\\n\"")

		// if i == len(lines)-1 {
		// 	result.WriteString(";")
		// }
		if i != len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
