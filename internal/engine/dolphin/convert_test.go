package dolphin_test

import (
	"strings"
	"testing"

	"github.com/kyleconroy/sqlc/internal/engine/dolphin"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
)

func Test_ParamRefParsing(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		paramRefCount int
	}{
		{name: "tuple comparison", paramRefCount: 2, input: "SELECT * WHERE (a, b) > (?, ?)"},
		{name: "cast", paramRefCount: 1, input: "SELECT CAST(? AS JSON)"},
		{name: "convert", paramRefCount: 1, input: "SELECT CONVERT(? USING UTF8)"},
		{name: "issues/1622", paramRefCount: 1, input: "INSERT INTO foo (x) VALUES (CAST(CONVERT(? USING UTF8) AS JSON))"},
	}
	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := dolphin.NewParser()
			stmts, err := p.Parse(strings.NewReader(test.input))
			if err != nil {
				t.Fatal(err)
			}
			if l := len(stmts); l != 1 {
				t.Fatalf("expected 1 statement, got %d", l)
			}
			paramRefs := extractParamRefs(stmts[0].Raw.Stmt)
			if got, want := len(paramRefs), test.paramRefCount; got != want {
				t.Fatalf("extracted params: want %d, got %d", want, got)
			}
		})
	}
}

// extractParamRefs extracts all of the ParamRef instances by walking the provided AST.
func extractParamRefs(n ast.Node) []*ast.ParamRef {
	var params []*ast.ParamRef
	astutils.Walk(astutils.VisitorFunc(func(n ast.Node) {
		switch t := n.(type) {
		case *ast.ParamRef:
			params = append(params, t)
		}
	}), n)
	return params
}

func Test_ParsingDoesNotPanic(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{name: "update implicit join", input: "UPDATE foo, bar SET foo.a = 1, bar.b = 2"},
	}
	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := dolphin.NewParser()
			_, err := p.Parse(strings.NewReader(test.input))
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
