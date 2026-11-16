// Package listing selects pages of a view's rows: search, filters, sort,
// limit and offset. models' views generator writes a List<Type> method per
// view that takes a Query; an app, or the package that renders its admin
// tables, fills the Query. The contract is models' interface §4.4.
package listing

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Query selects a page of a view's rows. Column names are SQL names.
type Query struct {
	Search   string   // rows where a SearchIn column, as text, contains it, case ignored
	SearchIn []string // column names; needed when Search isn't ""
	Filters  []Filter // every one must hold
	Sort     []Sort   // in order
	Limit    int      // 0: no limit
	Offset   int
}

// Filter keeps the rows whose Column compares to Value by Op.
type Filter struct {
	Column string
	Op     Op
	Value  any // passed to the driver, which encodes it by the column's type
}

// Op is a comparison of a Filter.
type Op string

// The comparisons a Filter can make.
const (
	Eq Op = "="
	Ne Op = "<>"
	Lt Op = "<"
	Le Op = "<="
	Gt Op = ">"
	Ge Op = ">="
)

// Sort orders the rows by Column, descending when Desc.
type Sort struct {
	Column string
	Desc   bool
}

// Source gives the query of a listing. Query implements it; so can a table
// package's own type. A generated List<Type> method takes a Source.
type Source interface {
	ListingQuery() Query
}

// ListingQuery returns q, so a Query is its own Source.
func (q Query) ListingQuery() Query { return q }

// Labeler gives the text for the id a link column is filtered by, such as
// a customer's name for its id, and false when there is no such row. A
// table calls it once per link filter of a request, to show the chosen
// value; rows carry their own text through the view's joins. Apps
// implement it; table packages call it. models generates none.
type Labeler interface {
	LabelOf(ctx context.Context, id string) (string, bool, error)
}

// ErrQuery is wrapped by every error of Build.
var ErrQuery = errors.New("listing: bad query")

func bad(format string, a ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{ErrQuery}, a...)...)
}

// Build returns the statement for a page of relation's rows and the
// statement counting the rows that match, with their arguments. For
// generated code. Both statements take the same arguments.
func Build(q Query, relation string, columns []string) (rows, count string, args []any, err error) {
	if relation == "" || strings.ContainsRune(relation, 0) {
		return "", "", nil, bad("relation %q", relation)
	}
	if len(columns) == 0 {
		return "", "", nil, bad("relation %q has no columns", relation)
	}
	known := make(map[string]bool, len(columns))
	for _, c := range columns {
		if c == "" || strings.ContainsRune(c, 0) || known[c] {
			return "", "", nil, bad("column %q of %q", c, relation)
		}
		known[c] = true
	}
	column := func(c string) (string, error) {
		if !known[c] {
			return "", bad("%q has no column %q", relation, c)
		}
		return ident(c), nil
	}
	if q.Limit < 0 {
		return "", "", nil, bad("negative limit %d", q.Limit)
	}
	if q.Offset < 0 {
		return "", "", nil, bad("negative offset %d", q.Offset)
	}

	var where []string
	searched := make([]string, 0, len(q.SearchIn))
	for _, c := range q.SearchIn {
		id, err := column(c)
		if err != nil {
			return "", "", nil, err
		}
		searched = append(searched, id)
	}
	if q.Search != "" {
		if strings.ContainsRune(q.Search, 0) {
			return "", "", nil, bad("search holds a NUL")
		}
		if len(searched) == 0 {
			return "", "", nil, bad("search %q with no column to search in", q.Search)
		}
		args = append(args, "%"+escapeLike(q.Search)+"%")
		n := "$" + strconv.Itoa(len(args))
		ors := make([]string, len(searched))
		for i, id := range searched {
			ors[i] = "CAST(" + id + " AS text) ILIKE " + n + ` ESCAPE '\'`
		}
		where = append(where, "("+strings.Join(ors, " OR ")+")")
	}
	for _, f := range q.Filters {
		id, err := column(f.Column)
		if err != nil {
			return "", "", nil, err
		}
		switch f.Op {
		case Eq, Ne, Lt, Le, Gt, Ge:
		default:
			return "", "", nil, bad("filter on %q: unknown op %q", f.Column, string(f.Op))
		}
		if isNil(f.Value) {
			return "", "", nil, bad("filter on %q: nil value", f.Column)
		}
		args = append(args, f.Value)
		where = append(where, id+" "+string(f.Op)+" $"+strconv.Itoa(len(args)))
	}
	var order []string
	for _, s := range q.Sort {
		id, err := column(s.Column)
		if err != nil {
			return "", "", nil, err
		}
		if s.Desc {
			id += " DESC"
		}
		order = append(order, id)
	}

	from := " FROM " + ident(relation)
	if len(where) > 0 {
		from += " WHERE " + strings.Join(where, " AND ")
	}
	selected := make([]string, len(columns))
	for i, c := range columns {
		selected[i] = ident(c)
	}
	var b strings.Builder
	b.WriteString("SELECT " + strings.Join(selected, ", ") + from)
	if len(order) > 0 {
		b.WriteString(" ORDER BY " + strings.Join(order, ", "))
	}
	if q.Limit > 0 {
		b.WriteString(" LIMIT " + strconv.Itoa(q.Limit))
	}
	if q.Offset > 0 {
		b.WriteString(" OFFSET " + strconv.Itoa(q.Offset))
	}
	return b.String(), "SELECT count(*)" + from, args, nil
}

// isNil reports a nil Value: untyped, or a nil pointer, map, slice,
// channel, function or interface, which drivers encode as NULL.
func isNil(v any) bool {
	if v == nil {
		return true
	}
	switch r := reflect.ValueOf(v); r.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return r.IsNil()
	}
	return false
}

// ident quotes a PostgreSQL identifier.
func ident(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// escapeLike escapes LIKE's wildcards and its escape character, \, in one
// pass, so no escape is escaped twice.
func escapeLike(s string) string { return likeEscaper.Replace(s) }
