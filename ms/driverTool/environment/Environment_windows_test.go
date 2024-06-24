package environment

import "testing"

func TestName(t *testing.T) {
	p := New()
	p.WalkDirs()
	p.Orig()
	p.Update()
}
