package gentest

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/epistemer/genx/genxtest"
)

var update = flag.Bool("update", false, "rewrite the golden files")

// root is models' module root, which holds the metafile.
const root = ".."

const domain = "sqlc.request/v1"

var shopRequest = genxtest.Input{Path: "../testdata/shop/request.json", Domain: domain}

// render runs one generator twice and fails unless both runs give the same
// bytes (roadmap, definition of done 2).
func render(t *testing.T, generator string, params map[string]any, inputs ...genxtest.Input) map[string][]byte {
	t.Helper()
	first := genxtest.RenderParams(t, root, generator, params, inputs...)
	second := genxtest.RenderParams(t, root, generator, params, inputs...)
	if len(first) != len(second) {
		t.Fatalf("%s: two runs wrote %v and %v", generator, names(first), names(second))
	}
	for name, b := range first {
		if !bytes.Equal(b, second[name]) {
			t.Fatalf("%s: two runs gave different bytes for %s", generator, name)
		}
	}
	return first
}

func names(files map[string][]byte) []string {
	var ns []string
	for n := range files {
		ns = append(ns, n)
	}
	slices.Sort(ns)
	return ns
}

// only returns the one file a render wrote, failing unless it is name.
func only(t *testing.T, files map[string][]byte, name string) []byte {
	t.Helper()
	if got := names(files); len(got) != 1 || got[0] != name {
		t.Fatalf("wrote %v, want [%s]", got, name)
	}
	return files[name]
}

func golden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", "golden", name)
	if *update {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%s differs from its golden file; run go test -update and review the diff.\ngot:\n%s", name, got)
	}
}

// TestDefaults renders the shop fixture with every parameter at its
// default: package models in models/, no prefix. (genx.yaml covers the
// shop with its prefix, through TestGenerated.)
func TestDefaults(t *testing.T) {
	golden(t, "shop-defaults.labels.json",
		only(t, render(t, "labels", nil, shopRequest), "models/zz_models.labels.json"))
	golden(t, "shop-defaults_enums.go.txt",
		only(t, render(t, "enums", nil, shopRequest), "models/zz_models_enums.go"))
	golden(t, "shop-defaults_views.go.txt",
		only(t, render(t, "views", nil, shopRequest), "models/zz_models_views.go"))
}

// TestDirAndPackage: the files follow dir and package.
func TestDirAndPackage(t *testing.T) {
	params := map[string]any{"dir": "internal/store", "package": "db"}
	labels := only(t, render(t, "labels", params, shopRequest), "internal/store/zz_db.labels.json")
	for _, want := range []string{`"dir": "internal/store"`, `"name": "db"`} {
		if !bytes.Contains(labels, []byte(want)) {
			t.Errorf("labels file lacks %s", want)
		}
	}
	enums := only(t, render(t, "enums", params, shopRequest), "internal/store/zz_db_enums.go")
	if !bytes.Contains(enums, []byte("\npackage db\n")) {
		t.Errorf("enums file isn't in package db:\n%s", enums)
	}
	views := only(t, render(t, "views", params, shopRequest), "internal/store/zz_db_views.go")
	if !bytes.Contains(views, []byte("\npackage db\n")) {
		t.Errorf("views file isn't in package db:\n%s", views)
	}
}

// TestPlain: a schema with no enums gets labels and no enums file.
func TestPlain(t *testing.T) {
	plain := genxtest.Input{Path: "../testdata/plain/request.json", Domain: domain}
	golden(t, "plain.labels.json",
		only(t, render(t, "labels", nil, plain), "models/zz_models.labels.json"))
	if files := render(t, "enums", nil, plain); len(files) != 0 {
		t.Errorf("enums wrote %v for a schema without enums", names(files))
	}
	if files := render(t, "views", nil, plain); len(files) != 0 {
		t.Errorf("views wrote %v for a schema without _view relations", names(files))
	}
}

// request writes a small request with the given schemas and queries, as
// JSON fragments, and returns it as an input.
func request(t *testing.T, defaultSchema, schemas, queries string) genxtest.Input {
	t.Helper()
	return engineRequest(t, "postgresql", defaultSchema, schemas, queries)
}

// engineRequest is request for another engine.
func engineRequest(t *testing.T, engine, defaultSchema, schemas, queries string) genxtest.Input {
	t.Helper()
	body := fmt.Sprintf(`{"settings": {"engine": %q}, "catalog": {"default_schema": %q, "schemas": [%s]}, "queries": [%s]}`,
		engine, defaultSchema, schemas, queries)
	path := filepath.Join(t.TempDir(), "request.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return genxtest.Input{Path: path, Domain: domain}
}

func public(tables, enums string) string {
	return fmt.Sprintf(`{"name": "public", "tables": [%s], "enums": [%s]}`, tables, enums)
}

func table(name string, columns ...string) string {
	var cs []string
	for _, c := range columns {
		cs = append(cs, fmt.Sprintf(`{"name": %q, "comment": ""}`, c))
	}
	return fmt.Sprintf(`{"rel": {"name": %q}, "comment": "", "columns": [%s]}`, name, strings.Join(cs, ", "))
}

func enum(name string, values ...string) string {
	var vs []string
	for _, v := range values {
		vs = append(vs, fmt.Sprintf("%q", v))
	}
	return fmt.Sprintf(`{"name": %q, "vals": [%s]}`, name, strings.Join(vs, ", "))
}

// TestNothingToWrite: no request, or a schema that gives no label, writes
// nothing; an enum with no values still gets its Label().
func TestNothingToWrite(t *testing.T) {
	for _, g := range []string{"labels", "enums", "views"} {
		if files := render(t, g, nil); len(files) != 0 {
			t.Errorf("%s with no request wrote %v", g, names(files))
		}
		if files := render(t, g, nil, request(t, "public", public("", ""), "")); len(files) != 0 {
			t.Errorf("%s with an empty schema wrote %v", g, names(files))
		}
	}
	empty := request(t, "public", public("", enum("nothing")), "")
	if files := render(t, "labels", nil, empty); len(files) != 0 {
		t.Errorf("labels for an enum with no values wrote %v", names(files))
	}
	got := only(t, render(t, "enums", nil, empty), "models/zz_models_enums.go")
	if !bytes.Contains(got, []byte("func (t Nothing) Label() i18n.TKer {")) {
		t.Errorf("no Label() for an enum with no values:\n%s", got)
	}
}

// TestOtherSystemSchemas: tables in pg_catalog, information_schema and
// pg_temp are skipped, not refused.
func TestOtherSystemSchemas(t *testing.T) {
	for _, schema := range []string{"pg_catalog", "information_schema", "pg_temp"} {
		in := request(t, "public",
			public(table("note", "id"), "")+`, {"name": "`+schema+`", "tables": [`+table("sys_table", "oid")+`], "enums": [`+enum("sys_enum", "a")+`]}`, "")
		labels := only(t, render(t, "labels", nil, in), "models/zz_models.labels.json")
		if bytes.Contains(labels, []byte("sys_")) {
			t.Errorf("labels for %s:\n%s", schema, labels)
		}
	}
}

// TestPackageForADirThatIsNotOne: package names the package when dir's last
// element couldn't.
func TestPackageForADirThatIsNotOne(t *testing.T) {
	params := map[string]any{"dir": "my-models", "package": "models"}
	only(t, render(t, "labels", params, shopRequest), "my-models/zz_models.labels.json")
	only(t, render(t, "enums", params, shopRequest), "my-models/zz_models_enums.go")
}

// TestNotReservedOutsideASCII: sqlc keeps letters outside ASCII in its Go
// names, so a table "él10n" is "Él10n" in sqlc, not reserved.
func TestNotReservedOutsideASCII(t *testing.T) {
	in := request(t, "public", public(table("él10n", "id"), ""), "")
	only(t, render(t, "labels", nil, in), "models/zz_models.labels.json")
}

// TestParamsWithoutRequest: parameters are checked on every run (§3).
func TestParamsWithoutRequest(t *testing.T) {
	for _, params := range []map[string]any{{"dir": "models/"}, {"package": "type"}, {"dir": "my models", "package": "models"}} {
		for _, g := range []string{"labels", "enums", "views"} {
			f := &failure{TB: t}
			genxtest.RenderParams(f, root, g, params)
			if !strings.Contains(f.msg, "models: parameter") {
				t.Errorf("%s with %v and no request: message %q, want a parameter refusal", g, params, f.msg)
			}
		}
	}
}

// failure records the message genxtest would fail the test with.
type failure struct {
	testing.TB
	msg string
}

func (f *failure) Helper() {}
func (f *failure) Fatalf(format string, args ...any) {
	if f.msg == "" {
		f.msg = fmt.Sprintf(format, args...)
	}
}

// TestRefusals: every case of interface §6 shared by the generators, for
// each of them.
func TestRefusals(t *testing.T) {
	note := func(t *testing.T) []genxtest.Input {
		return []genxtest.Input{request(t, "public", public(table("note", "id"), enum("kind", "a")), "")}
	}
	for _, c := range []struct {
		name   string
		params map[string]any
		inputs func(t *testing.T) []genxtest.Input
		want   string
	}{
		{"dir with a trailing slash", map[string]any{"dir": "models/"}, note, `models: parameter dir "models/": not a clean relative slash-separated path`},
		{"dir outside", map[string]any{"dir": "../models"}, note, `models: parameter dir "../models": not a clean relative slash-separated path`},
		{"dir absolute", map[string]any{"dir": "/models"}, note, `models: parameter dir "/models": not a clean relative slash-separated path`},
		{"dir empty element", map[string]any{"dir": "a//b"}, note, `models: parameter dir "a//b": not a clean relative slash-separated path`},
		{"dir dot element", map[string]any{"dir": "a/./b"}, note, `models: parameter dir "a/./b": not a clean relative slash-separated path`},
		{"dir empty", map[string]any{"dir": ""}, note, `models: parameter dir "": not a clean relative slash-separated path`},
		{"dir dot", map[string]any{"dir": "."}, note, `models: parameter dir ".": not a clean relative slash-separated path`},
		{"dir with a space", map[string]any{"dir": "my models", "package": "models"}, note, `models: parameter dir "my models": not a clean relative slash-separated path`},
		{"dir hidden", map[string]any{"dir": ".models", "package": "models"}, note, `models: parameter dir ".models": not a clean relative slash-separated path`},
		{"dir reserved on Windows", map[string]any{"dir": "db/aux.d", "package": "models"}, note, `models: parameter dir "db/aux.d": not a clean relative slash-separated path`},
		{"package empty", map[string]any{"package": ""}, note, `models: parameter package: "" is not a Go package name`},
		{"dir not a package name", map[string]any{"dir": "my-models"}, note, `models: parameter package: "my-models", the last element of dir, is not a Go package name; set package`},
		{"package with a dash", map[string]any{"package": "my-models"}, note, `models: parameter package: "my-models" is not a Go package name`},
		{"package keyword", map[string]any{"package": "type"}, note, `models: parameter package: "type" is not a Go package name`},
		{"package blank", map[string]any{"package": "_"}, note, `models: parameter package: "_" is not a Go package name`},
		{"package number", map[string]any{"package": 7}, note, `models: parameter package: 7 is not a Go package name`},
		{"two requests", nil, func(t *testing.T) []genxtest.Input {
			return append(note(t), note(t)...)
		}, `models: 2 sqlc requests; one is supported`},
		{"default schema missing", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "main", public(table("note", "id"), ""), "")}
		}, `models: the default schema "main" isn't in the request`},
		{"another schema with only an enum", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public",
				public(table("note", "id"), "")+`, {"name": "audit", "tables": [], "enums": [`+enum("level", "a")+`]}`, "")}
		}, `models: schema "audit" has tables or enums; only the default schema "public" is supported`},
		{"another schema", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public",
				public(table("note", "id"), "")+`, {"name": "audit", "tables": [`+table("log", "id")+`], "enums": []}`, "")}
		}, `models: schema "audit" has tables or enums; only the default schema "public" is supported`},
		{"enum outside ASCII", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public("", enum("état", "a")), "")}
		}, `models: enum "état": sqlc's Go name for it can't be derived`},
		{"enum with no Go name", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public("", enum("_", "a")), "")}
		}, `models: enum "_": sqlc's Go name for it can't be derived`},
		{"table not an identifier", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public(table("first name", "id"), ""), "")}
		}, `models: table "first name" gives "First name", which isn't a Go identifier`},
		{"table empty after the prefix", map[string]any{"table_prefix": "app_"}, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public(table("app_", "id"), ""), "")}
		}, `models: table "app_" gives "", which isn't a Go identifier`},
		{"column with no name", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public(table("note", "_"), ""), "")}
		}, `models: column "note"."_" gives "", which isn't a Go identifier`},
		{"two columns, one key", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public(table("note", "user1", "user_1"), ""), "")}
		}, `models: column "note"."user1" and column "note"."user_1" both give the key "table.note.Note_User1"`},
		{"two tables, one variable", map[string]any{"table_prefix": "app_"}, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public(table("app_user")+", "+table("user"), ""), "")}
		}, `models: table "app_user" and table "user" both give the variable "L10nUser"`},
		{"sqlc's numbering clash", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public("", enum("e", "value_1", "!!")), "")}
		}, `models: enum "e" value "value_1" and enum "e" value "!!" both give the variable "L10nEValue1"`},
		{"two enums, one Go type", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public("", enum("user1", "a")+", "+enum("user_1", "b")), "")}
		}, `models: enum "user1" and enum "user_1" both give sqlc the Go type "User1"`},
		{"an option list named like an enum", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public("", enum("x", "a")+", "+enum("all_x_options", "b")), "")}
		}, `models: enum "x" gives the function "AllXOptions", which is enum "all_x_options"'s Go type`},
		{"an option list named like an enum constant", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public("", enum("x", "a")+", "+enum("all_x", "options")), "")}
		}, `models: enum "x" gives the function "AllXOptions", which is a constant of enum "all_x"`},
		{"table reserved", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public(table("l10n_note", "id"), ""), "")}
		}, `models: table "l10n_note" gives sqlc the Go name "L10nNote"; names starting with "L10n" are reserved for labels`},
		{"enum reserved", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public("", enum("l10n_kind", "a")), "")}
		}, `models: enum "l10n_kind" gives sqlc the Go name "L10nKind"; names starting with "L10n" are reserved for labels`},
		{"query reserved", nil, func(t *testing.T) []genxtest.Input {
			return []genxtest.Input{request(t, "public", public(table("note", "id"), ""), `{"name": "L10nNotes"}`)}
		}, `models: query "L10nNotes" gives sqlc the Go name "L10nNotes"; names starting with "L10n" are reserved for labels`},
	} {
		for _, g := range []string{"labels", "enums", "views"} {
			t.Run(c.name+"/"+g, func(t *testing.T) {
				f := &failure{TB: t}
				files := genxtest.RenderParams(f, root, g, c.params, c.inputs(t)...)
				if f.msg == "" {
					t.Fatalf("no refusal; wrote %v", names(files))
				}
				if !strings.Contains(f.msg, c.want) {
					t.Errorf("message:\n%s\nwant it to contain:\n%s", f.msg, c.want)
				}
			})
		}
	}
}

// TestViewsNames: sqlc's names for a digit-led relation, an oxen relation
// (sqlc's one singular rule that ignores the end) and a keyword column; a
// query named Listings declares no type, so it isn't refused; a relation
// whose name merely contains _view gets nothing.
func TestViewsNames(t *testing.T) {
	in := request(t, "public", public(table("1_view", "order", "user_id")+", "+table("Oxen_view", "id")+", "+
		table("view_log", "id")+", "+table("preview", "id"), ""), `{"name": "Listings"}`)
	got := string(only(t, render(t, "views", nil, in), "models/zz_models_views.go"))
	for _, want := range []string{
		"func (q *Queries) List_1View(ctx context.Context, src listing.Source) ([]_1View, int64, error) {",
		`listing.Build(lq, "1_view", []string{`,
		"&i.Order,", "&i.UserID,",
		"func (q *Queries) ListOxView(ctx context.Context, src listing.Source) ([]OxView, int64, error) {",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("views file lacks %s:\n%s", want, got)
		}
	}
	for _, unwanted := range []string{"ViewLog", "Preview"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("views file has %s:\n%s", unwanted, got)
		}
	}
}

// TestViewsOtherEngine: another engine is refused only when there is a
// _view relation to list.
func TestViewsOtherEngine(t *testing.T) {
	in := engineRequest(t, "mysql", "public", public(table("note", "id"), ""), "")
	if files := render(t, "views", nil, in); len(files) != 0 {
		t.Errorf("views wrote %v", names(files))
	}
}

// TestViewsRefusals: the cases of interface §6 that are views' own. labels
// and enums make none of them.
func TestViewsRefusals(t *testing.T) {
	in := func(tables, enums, queries string) func(t *testing.T) genxtest.Input {
		return func(t *testing.T) genxtest.Input { return request(t, "public", public(tables, enums), queries) }
	}
	for _, c := range []struct {
		name  string
		input func(t *testing.T) genxtest.Input
		want  string
	}{
		{"another engine", func(t *testing.T) genxtest.Input {
			return engineRequest(t, "mysql", "public", public(table("order_view", "id"), ""), "")
		}, `models: views need PostgreSQL; the request's engine is "mysql"`},
		{"view outside ASCII", in(table("café_view", "id"), "", ""), `models: view "café_view": sqlc's Go name for it can't be derived`},
		{"column outside ASCII", in(table("order_view", "prénom"), "", ""), `models: view "order_view" column "prénom": sqlc's Go name for it can't be derived`},
		{"view with no columns", in(table("empty_view"), "", ""), `models: view "empty_view" has no columns`},
		// zp's casing keeps these apart (Id, ID), so only views sees them.
		{"two columns, one field", in(table("dup_view", "id", "ID"), "", ""), `models: view "dup_view": columns "id" and "ID" both give the field "ID"`},
		{"two views, one type", in(table("id_view", "x")+", "+table("ID_view", "x"), "", ""), `models: views "id_view" and "ID_view" both give sqlc the Go type "IDView"`},
		{"query named like the method", in(table("order_view", "id"), "", `{"name": "ListOrderView"}`), `models: query "ListOrderView" clashes with the listing method of view "order_view"`},
		{"query with the constant listing", in(table("order_view", "id"), "", `{"name": "Listing"}`), `models: query "Listing" gives sqlc the constant "listing", which hides a name the views file uses`},
		{"query named listing", in(table("order_view", "id"), "", `{"name": "listing"}`), `models: query "listing" gives sqlc the constant "listing", which hides a name the views file uses`},
		{"query hiding len", in(table("order_view", "id"), "", `{"name": "Len"}`), `models: query "Len" gives sqlc the constant "len", which hides a name the views file uses`},
		{"query hiding column", in(table("order_view", "id"), "", `{"name": "Column"}`), `models: query "Column" gives sqlc the constant "column", which hides a name the views file uses`},
		{"query named like a column set", in(table("order_view", "id"), "", `{"name": "OrderViewColumns"}`), `models: view "order_view": its column set "OrderViewColumns" is already a query's, a relation's or an enum's Go name`},
		{"relation named like a column set", in(table("order_view", "id")+", "+table("order_view_columns", "id"), "", ""), `models: view "order_view": its column set "OrderViewColumns" is already a query's, a relation's or an enum's Go name`},
		{"query hiding context", in(table("order_view", "id"), "", `{"name": "Context"}`), `models: query "Context" gives sqlc the constant "context", which hides a name the views file uses`},
		{"query named like an oxen method", in(table("oxen_view", "id"), "", `{"name": "ListOxView"}`), `models: query "ListOxView" clashes with the listing method of view "oxen_view"`},
		{"enum named Listings", in(table("order_view", "id"), enum("listings", "a"), ""), `models: enum "listings" gives sqlc the Go name "Listings", which views declares`},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := &failure{TB: t}
			files := genxtest.RenderParams(f, root, "views", nil, c.input(t))
			if f.msg == "" {
				t.Fatalf("no refusal; wrote %v", names(files))
			}
			if !strings.Contains(f.msg, c.want) {
				t.Errorf("message:\n%s\nwant it to contain:\n%s", f.msg, c.want)
			}
			for _, g := range []string{"labels", "enums"} {
				f := &failure{TB: t}
				genxtest.RenderParams(f, root, g, nil, c.input(t))
				if f.msg != "" {
					t.Errorf("%s refused it too: %s", g, f.msg)
				}
			}
		})
	}
}

// TestViewsControlCharacter: a newline in a view's name never reaches the
// generated comment; the shared rule refuses it first.
func TestViewsControlCharacter(t *testing.T) {
	f := &failure{TB: t}
	genxtest.RenderParams(f, root, "views", nil, request(t, "public", public(table("order\n_view", "id"), ""), ""))
	if !strings.Contains(f.msg, "isn't a Go identifier") {
		t.Errorf("message %q", f.msg)
	}
}
