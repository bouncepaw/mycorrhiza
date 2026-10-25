package misc

import (
	"reflect"
	"testing"
)

func TestNestList(t *testing.T) {
	t.Parallel()

	flat := []listDatum{
		{Name: "alpha"},
		{Name: "alpha/child", Ext: "png"},
		{Name: "alpha/missing/leaf"},
		{Name: "alphabet"},
	}
	wantTree := []*listDatum{
		{
			Name: "alpha", Label: "alpha", Exists: true,
			Children: []*listDatum{
				{Name: "alpha/child", Label: "child", Ext: "png", Exists: true},
				{
					Name: "alpha/missing", Label: "missing",
					Children: []*listDatum{
						{Name: "alpha/missing/leaf", Label: "leaf", Exists: true},
					},
				},
			},
		},
		{Name: "alphabet", Label: "alphabet", Exists: true},
	}
	if got := nestList(flat); !reflect.DeepEqual(got, wantTree) {
		t.Errorf("nestList returned an unexpected tree: got %+v, want %+v", got, wantTree)
	}
	if got := nestList(nil); len(got) != 0 {
		t.Errorf("nestList(nil) = %v, want an empty list", got)
	}
}
