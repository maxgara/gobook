package live

import (
	"strings"
	"testing"
)

func TestLive(t *testing.T) {
	var x Xtype
	LivePrint(&x)
}

type Xtype struct {
	val string
	fc  *Xtype
	ns  *Xtype
}

func (x *Xtype) WString() string {
	return x.val
}
func (x *Xtype) WInput(input string) {
	dstrs := strings.Split(x.val, input)
	x.fc = &Xtype{val: dstrs[0]}
	prior := x.fc
	if len(dstrs) == 1 {
		return
	}
	for _, v := range dstrs[1:] {
		newnode := &Xtype{val: v}
		(*prior).ns = newnode
		prior = newnode
	}
}
func (x *Xtype) WInit() {
	*x = Xtype{val: "this is a test\ndon't panic."}
}

func (x *Xtype) FirstChild() bool {
	if x.fc == nil {
		return false
	}
	*x = *x.fc
	return true
}
func (x *Xtype) NextSibling() bool {
	if x.ns == nil {
		return false
	}
	*x = *x.ns
	return true
}
