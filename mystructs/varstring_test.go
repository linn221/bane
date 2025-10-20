package mystructs

import "testing"

func TestVarString_ParseEqualsPlaceholders(t *testing.T) {
	cases := []struct {
		in       string
		wantExec string
		inject   map[string]string
	}{
		{"hello-{name=world}", "hello-world", nil},
		{"{id=1}-{name=Henry Cohle}", "1-Henry Cohle", nil},
		{"{id=1}-{name=Henry Cohle}", "2-Henry Cohle", map[string]string{"id": "2"}},
		{"pre-{a=}-{b=ok}-post", "pre--ok-post", nil},
		{"x{n=0}y{m=1}z", "x0y1z", nil},
	}
	for i, c := range cases {
		vs, err := NewVarString(c.in)
		if err != nil {
			t.Fatalf("case %d: unexpected err: %v", i, err)
		}
		if c.inject != nil {
			vs.Inject(c.inject)
		}
		got := vs.Exec()
		if got != c.wantExec {
			t.Errorf("case %d: Exec()=%q want %q", i, got, c.wantExec)
		}
	}
}
