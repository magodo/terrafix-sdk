package hclutils_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/magodo/terrafix-sdk/internal/hclutils"
	"github.com/zclconf/go-cty/cty"
)

func TestTraversalMatches(t *testing.T) {
	cases := []struct {
		name string
		t1   string
		t2   string
		res  bool
		err  bool
	}{
		{
			name: "emty traversal unmatch",
			t1:   "",
			t2:   "a.b",
			res:  false,
		},
		{
			name: "emty sub-traversal failed to parse",
			t1:   "a.b",
			t2:   "",
			err:  true,
		},
		{
			name: "regular match 1",
			t1:   "a.0.b",
			t2:   "a.b",
			res:  true,
		},
		{
			name: "regular match 2",
			t1:   "a[0].b",
			t2:   "a.b",
			res:  true,
		},
		{
			name: "regular match 3",
			t1:   "a[0].b.1",
			t2:   "a.b.1",
			res:  true,
		},
		{
			name: "regular unmatch 1",
			t1:   "a[0].b.1",
			t2:   "b",
			res:  false,
		},
		{
			name: "regular unmatch 2",
			t1:   "a[0].b.1",
			t2:   "a.b",
			res:  false,
		},
		{
			name: "regular unmatch 3",
			t1:   "a[0].b.1.c",
			t2:   "a.b",
			res:  false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var t1 hcl.Traversal
			if tt.t1 != "" {
				var err error
				t1, err = hclutils.ParseTraversal(tt.t1)
				if err != nil {
					t.Fatal(err.Error())
				}
			}
			got, err := hclutils.TraversalMatches(t1, tt.t2)
			if tt.err {
				if err == nil {
					t.Fatal("expect error, but none")
				}
				return
			}
			if err != nil {
				t.Fatalf(err.Error())
			}
			if got != tt.res {
				t.Fatalf("find address in traversal failed: expect=%t, got=%t", tt.res, got)
			}
		})
	}
}

func TestFindAddrInTraversal(t *testing.T) {
	cases := []struct {
		name string
		t1   string
		t2   string
		idx  int
		err  bool
	}{
		{
			name: "emty traversal unmatch",
			t1:   "",
			t2:   "a.b",
			idx:  -1,
		},
		{
			name: "emty sub-traversal failed to parse",
			t1:   "a.b",
			t2:   "",
			err:  true,
		},
		{
			name: "regular match 1",
			t1:   "a.0.b.1",
			t2:   "a.b",
			idx:  2,
		},
		{
			name: "regular match 2",
			t1:   "a[0].b.1",
			t2:   "a.b",
			idx:  2,
		},
		{
			name: "regular match 3",
			t1:   "a[0].b.1",
			t2:   "a.b.1",
			idx:  3,
		},
		{
			name: "regular match 4",
			t1:   "0.a.b.1",
			t2:   "a.b",
			idx:  2,
		},
		{
			name: "regular unmatch 1",
			t1:   "a[0].b.1",
			t2:   "b",
			idx:  -1,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var t1 hcl.Traversal
			if tt.t1 != "" {
				var err error
				// The "foo" prefix here is to make the address parsing happy
				// to ensure the first step is not a number.
				t1, err = hclutils.ParseTraversal("foo." + tt.t1)
				if err != nil {
					t.Fatal(err.Error())
				}
				t1 = t1[1:]
			}
			got, err := hclutils.FindSubAddr(t1, tt.t2)
			if tt.err {
				if err == nil {
					t.Fatal("expect error, but none")
				}
				return
			}
			if err != nil {
				t.Fatalf(err.Error())
			}
			if got != tt.idx {
				t.Fatalf("find address in traversal failed: expect=%d, got=%d", tt.idx, got)
			}
		})
	}
}

func TestFormatTraversal(t *testing.T) {
	cases := []struct {
		name   string
		input  hcl.Traversal
		output string
	}{
		{
			name:   "empty",
			input:  nil,
			output: "",
		},
		{
			name: `a.b[0].c[key].d[*]`,
			input: hcl.Traversal{
				hcl.TraverseRoot{Name: "a"},
				hcl.TraverseAttr{Name: "b"},
				hcl.TraverseIndex{Key: cty.NumberIntVal(0)},
				hcl.TraverseAttr{Name: "c"},
				hcl.TraverseIndex{Key: cty.StringVal("key")},
				hcl.TraverseAttr{Name: "d"},
				hcl.TraverseSplat{},
			},
			output: `a.b[0].c[key].d[*]`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if actual := hclutils.FormatTraversal(tt.input); tt.output != actual {
				t.Fatalf("expect=%s, got=%s", tt.output, actual)
			}
		})
	}
}
