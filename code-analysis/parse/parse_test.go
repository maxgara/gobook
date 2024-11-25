package parse

import (
	"fmt"
	"strings"
)

const test = `func f1(c int, d float) {
	var x int
	var y int
	var z int
	z = x / 2
	z = z * 2
	y = z - 3
	fmt.Println(y)
}
	func f2(a, b string){
	return
	}`

//	func ExampleParseNd_Parse() {
//		s := "abc abc abd abx bay pav"
//		x := NewParseNd(s)
//		words := x.Parse("ab")
//		fmt.Print(words)
//	}
func ExampleParseNd_Parse() {
	s := test
	x := NewParseNd(s)
	_ = x.Parse("(?<myname>var .*)")
	fmt.Println(x)
	// Output:
}
func ExampleParseNd_Parse_Two() {
	s := "abc abd abx bay pab"
	x := NewParseNd(s)
	x.Temp("(?<myname>ab.)").Parse("(?<myname2>.)")
	fmt.Println(x)
	// Output:
}
func ExampleParseNd_Temp() {
	s := "abc abd abx bay pab"
	x := NewParseNd(s)
	x.Temp("(?<myname>ab.)").Temp("(?<myname2>.)").Parse("(?<myname3>.)")
	fmt.Println(x)
	//Output:
}
func ExamplePname() {
	p := "abc (?<myname>..)"
	pn := pname(p)
	fmt.Println(pn)
}
func ExamplePrintnd() {
	s := strings.ReplaceAll("abc abd abx bay pab", " ", "\n")
	x := NewParseNd(s)
	x.Parse(`(?<subnode1>a.)`)
	fmt.Println(x.stringnd(0, "mynode"))
}
func ExampleRstring() {
	s := strings.ReplaceAll("abc abd abx bay pab", " ", "\n")
	x := NewParseNd(s)
	_ = x.Parse(`(?<sub1grp>.\n.)`).Parse(`(?<sub2>.)`)
	fmt.Println(x.rString(0, "TEST"))
	fmt.Printf("%#v", x)
}
func ExamplePread() {
	s := strings.ReplaceAll("abc abd abx bay pab", " ", "\n")
	x := NewParseNd(s)
	_ = x.Parse(`(?<sub1grp>...)`).Parse(`(?<sub2>.)`)
	x.Parse(`(?<g2ds>\w\w)`)
	r := x.newPreader()
	for r.Read() != nil {
		fmt.Println(r.pinfo)
	}

}
func ExamplePText() {
	s := strings.ReplaceAll("abc abd abx bay pab", " ", "\n")
	x := NewParseNd(s)
	_ = x.Parse(`(?<sub1grp>...)`).Parse(`(?<sub2>.)`)
	x.Parse(`(?<g2ds>\w\w)`)
	r := x.newPreader()
	for r.Read() != nil {
		fmt.Println(r.Text())
	}
	//Output:
}
func ExampleParseNd_Name() {
	s := strings.ReplaceAll("abc abd abx bay pab", " ", "\n")
	x := NewParseNd(s)
	_ = x.Parse(`(?<sub1grp>...)`).Parse(`(?<sub2>.)`)
	q := x.p["sub1grp"][0]
	fmt.Printf("target node for .Name: %v\n", q)
	fmt.Printf("full parse tree: %v\n", x)
	fmt.Println("name: " + q.Name())
	//Output:
}
