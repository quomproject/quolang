package quosyntax

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"

	"github.com/hashicorp/hcl/v2"
)

func TestWalk(t *testing.T) {

	tests := []struct {
		src  string
		want []testWalkCall
	}{
		{
			`1`,
			[]testWalkCall{
				{testWalkEnter, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.LiteralValueExpr"},
			},
		},
		{
			`foo`,
			[]testWalkCall{
				{testWalkEnter, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.ScopeTraversalExpr"},
			},
		},
		{
			`1 + 1`,
			[]testWalkCall{
				{testWalkEnter, "*quosyntax.BinaryOpExpr"},
				{testWalkEnter, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.LiteralValueExpr"},
				{testWalkEnter, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.BinaryOpExpr"},
			},
		},
		{
			`(1 + 1)`,
			[]testWalkCall{
				{testWalkEnter, "*quosyntax.BinaryOpExpr"},
				{testWalkEnter, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.LiteralValueExpr"},
				{testWalkEnter, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.BinaryOpExpr"},
			},
		},
		{
			`a[0]`,
			[]testWalkCall{
				// because the index is constant here, the index is absorbed into the traversal
				{testWalkEnter, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.ScopeTraversalExpr"},
			},
		},
		{
			`0[foo]`, // semantically incorrect, but should still parse and be walkable
			[]testWalkCall{
				{testWalkEnter, "*quosyntax.IndexExpr"},
				{testWalkEnter, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.LiteralValueExpr"},
				{testWalkEnter, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.IndexExpr"},
			},
		},
		{
			`bar()`,
			[]testWalkCall{
				{testWalkEnter, "*quosyntax.FunctionCallExpr"},
				{testWalkExit, "*quosyntax.FunctionCallExpr"},
			},
		},
		{
			`bar(1, a)`,
			[]testWalkCall{
				{testWalkEnter, "*quosyntax.FunctionCallExpr"},
				{testWalkEnter, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.LiteralValueExpr"},
				{testWalkEnter, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.FunctionCallExpr"},
			},
		},
		{
			`bar(1, a)[0]`,
			[]testWalkCall{
				{testWalkEnter, "*quosyntax.RelativeTraversalExpr"},
				{testWalkEnter, "*quosyntax.FunctionCallExpr"},
				{testWalkEnter, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.LiteralValueExpr"},
				{testWalkEnter, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.FunctionCallExpr"},
				{testWalkExit, "*quosyntax.RelativeTraversalExpr"},
			},
		},
		{
			`[for x in foo: x + 1 if x < 10]`,
			[]testWalkCall{
				{testWalkEnter, "*quosyntax.ForExpr"},
				{testWalkEnter, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.ScopeTraversalExpr"},
				{testWalkEnter, "quosyntax.ChildScope"},
				{testWalkEnter, "*quosyntax.BinaryOpExpr"},
				{testWalkEnter, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.ScopeTraversalExpr"},
				{testWalkEnter, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.BinaryOpExpr"},
				{testWalkExit, "quosyntax.ChildScope"},
				{testWalkEnter, "quosyntax.ChildScope"},
				{testWalkEnter, "*quosyntax.BinaryOpExpr"},
				{testWalkEnter, "*quosyntax.ScopeTraversalExpr"},
				{testWalkExit, "*quosyntax.ScopeTraversalExpr"},
				{testWalkEnter, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.LiteralValueExpr"},
				{testWalkExit, "*quosyntax.BinaryOpExpr"},
				{testWalkExit, "quosyntax.ChildScope"},
				{testWalkExit, "*quosyntax.ForExpr"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
			expr, diags := ParseExpression([]byte(test.src), "", hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				t.Fatalf("failed to parse expression: %s", diags.Error())
			}

			w := testWalker{}
			diags = Walk(expr, &w)
			if diags.HasErrors() {
				t.Fatalf("failed to walk: %s", diags.Error())
			}

			got := w.Calls
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("wrong calls\ngot: %swant: %s", spew.Sdump(got), spew.Sdump(test.want))
				for _, problem := range deep.Equal(got, test.want) {
					t.Errorf(problem)
				}
			}
		})
	}
}

type testWalkMethod int

const testWalkEnter testWalkMethod = 1
const testWalkExit testWalkMethod = 2

type testWalkCall struct {
	Method   testWalkMethod
	NodeType string
}

type testWalker struct {
	Calls []testWalkCall
}

func (w *testWalker) Enter(node Node) hcl.Diagnostics {
	w.Calls = append(w.Calls, testWalkCall{testWalkEnter, fmt.Sprintf("%T", node)})
	return nil
}

func (w *testWalker) Exit(node Node) hcl.Diagnostics {
	w.Calls = append(w.Calls, testWalkCall{testWalkExit, fmt.Sprintf("%T", node)})
	return nil
}
