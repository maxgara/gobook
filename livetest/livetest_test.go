package live

import "testing"

func TestParseLive(t *testing.T) {
	var x Xtype
	LivePrint(&x)
}

type Xtype struct {
	s string
}

func (x *Xtype) WString() string {
	return x.s
}
func (x *Xtype) WInput(input string) {
	x.s += input
}
func (x *Xtype) WInit() {
	*x = Xtype{}
}
