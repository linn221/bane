package mystructs

import "testing"

func TestVarKVGroup_ParseColonPairs_WithVarStringValues(t *testing.T) {
	in := "header1:val1 Token:Bearer{token=abc123} UA:Mozilla{agentVersion=1.0}"
	var vkg VarKVGroup
	if err := vkg.UnmarshalGQL(in); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if len(vkg.VarKVs) != 3 {
		t.Fatalf("got %d pairs, want 3", len(vkg.VarKVs))
	}
	want := []struct{ k, v string }{
		{"header1", "val1"},
		{"Token", "Bearer{token=abc123}"},
		{"UA", "Mozilla{agentVersion=1.0}"},
	}
	for i, w := range want {
		if vkg.VarKVs[i].Key.OriginalString != w.k || vkg.VarKVs[i].Value.OriginalString != w.v {
			t.Errorf("pair %d got (%q,%q) want (%q,%q)", i, vkg.VarKVs[i].Key.OriginalString, vkg.VarKVs[i].Value.OriginalString, w.k, w.v)
		}
	}
}

func TestVarKVGroup_Exec_EvaluatesVarStrings(t *testing.T) {
	in := "A:{a=1} B:{b=} C:{c=3}"
	var vkg VarKVGroup
	if err := vkg.UnmarshalGQL(in); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	// Inject variables
	vkg.VarKVs[1].Value.Inject(map[string]string{"b": "2"})
	got := vkg.Exec()
	if got != "A:1 B:2 C:3" {
		t.Errorf("Exec=%q want %q", got, "A:1 B:2 C:3")
	}
}

func TestVarKVGroup_Roundtrip_Value_Scan(t *testing.T) {
	in := "k1:v1 k2:{v=2}"
	var vkg VarKVGroup
	if err := vkg.UnmarshalGQL(in); err != nil {
		t.Fatal(err)
	}
	dv, err := vkg.Value()
	if err != nil {
		t.Fatal(err)
	}
	var back VarKVGroup
	if err := back.Scan(dv); err != nil {
		t.Fatal(err)
	}
	if back.Exec() != "k1:v1 k2:2" {
		t.Errorf("roundtrip Exec=%q", back.Exec())
	}
}
