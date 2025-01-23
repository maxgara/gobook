package main

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

const str = `9 + 4 - 2`

func main() {
	fset := token.FileSet{}
	exp, err := parser.ParseExpr(str)
	if err != nil {
		fmt.Println(err)
		return
	}
	format.Node(os.Stdout, &fset, exp)
	//output: 9 + 4 - 2
}
