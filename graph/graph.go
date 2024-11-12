package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type dataline struct {
	nums []float64
	strs []string
	err  error
}
type datareader struct {
	lines []dataline
}

func (d *datareader) next() []float64 {
	//finish this
	return nil
}

func test() {
	fmt.Println(nextData())
}
func main() {
	test()
	os.Exit(0)
	//<svg viewBox="-10 -10 220 120" xmlns="http://www.w3.org/2000/svg">
	fmt.Printf(`<svg viewBox="-10 -10 220 120" xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>`, 600, 500)
	fmt.Printf(`<polyline stroke="black" fill="none" points="`)
	// "50,0 21,90 98,35 2,35 79,90" />
	p, head, err := nextData()
	var data [][]float64
	for {
		if head != nil {
			continue
		}
		//add index as x coord if data is single-column
		if len(p) == 1 {
			p = []float64{float64(len(data)), p[0]}
		}
		data = append(data, p)
		if len(data) != 1 {
			fmt.Printf(" ") //space before each new point
		}
		fmt.Printf("%f,%f", p[0], p[1])
		if err == io.EOF {
			break
		}
		p, head, _ = nextData()
	}
	fmt.Printf(` />`)
	// fmt.Printf("<polygon points='%.5g,%.5g %.5g,%.5g %.5g,%.5g %.5g,%.5g' stroke='#%s'/>\n")
}

func nextData() ([]float64, []string, error) {
	var buff = make([]byte, 100)
	n, err := os.Stdin.Read(buff)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}
	rows := strings.Split(string(buff[:n]), "\n")
	for _, row := range rows {
		cols := strings.Split(row, " ")
		fmt.Printf("cols:<<%v>>\n", cols)
		var data []float64
		var strs []string
		for _, d := range cols {
			fmt.Printf("convesion attempt on %v\n", d)
			//check if col is valid data
			x, err := strconv.ParseFloat(d, 64)
			//get header if not
			if err != nil {
				strs = append(strs, d)
				continue
			}
			data = append(data, x)
		}
	}

	return data, strs, err
}
