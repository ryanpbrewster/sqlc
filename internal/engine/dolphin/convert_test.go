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
			e := &refExtractor{}
			astutils.Walk(e, stmts[0].Raw.Stmt)
			if got, want := len(e.params), test.paramRefCount; got != want {
				t.Fatalf("extracted params: want %d, got %d", want, got)
			}
		})
	}
}

// refExtractor is an astutils.Visitor instance that will extract all of the
// ParamRef instances it encounters while walking an AST.
type refExtractor struct {
	params []*ast.ParamRef
}

func (e *refExtractor) Visit(n ast.Node) astutils.Visitor {
	switch t := n.(type) {
	case *ast.ParamRef:
		e.params = append(e.params, t)
	}
	return e
}
