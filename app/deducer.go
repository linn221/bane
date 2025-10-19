package app

import "sync"

type Deducer struct {
	ProgramId  int
	EndpointId int
	References []Reference
	mu         sync.Mutex
}

type Reference struct {
	ReferenceId   int
	ReferenceType string
}

// func (d *Deducer) SetReference(referenceId int, refType string) int {
// 	d.mu.Lock()
// 	defer d.mu.Unlock()
// 	d.References = append(d.References, Reference{ReferenceId: referenceId, ReferenceType: refType})
// 	return len(d.References) - 1
// }

func (d *Deducer) Lock() func() {
	d.mu.Lock()
	return d.mu.Unlock
}

func (d *Deducer) ReadRId(rid int) (int, string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ref := d.References[rid]
	return ref.ReferenceId, ref.ReferenceType
}
