package mystructs

type VarKVGroup struct {
	VarKVs []VarKV
}

type VarKV struct {
	Key   VarString
	Value VarString
}
