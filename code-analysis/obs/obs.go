package main

// object group
type Obs struct {
	arr []Ob //shared props
}

// object
type Ob struct {
	name string
	val  string
	p    map[string]Obs //props
}

// ex: Ob FileString -> Obg Function [3obs]-> Obg [fname] [1ob]
func (o *Ob) Parse(name, pattern string) Obs {
	strs := []string{"place", "holder"}
	if len(strs) == 0 {
		return Obs{}
	}
	// var g Obs = Obs{p: make(map[string]bool)}
	// for i, v := range strs {
	// 	newname := name + fmt.Sprintf("%v", i)
	// 	newob := Ob{name: newname, val: v, p: make(map[string]Obs)}
	// }
	// o.p[name] = g
	// Ob.p[name] = Obg
}

func (o Obs) ParseAll(name, pattern string) Obs {
	// TODO
	return Obs{}
}

// read prop p for all x in o and return set of all x.y.val
func (o Obs) Get(p string) []string {
	//TODO
	return nil
}
