package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestParse0(t *testing.T) {
	b := bytes.NewBuffer([]byte(data))
	p := parser{r: b}
	fmt.Println(p.parse())
	fmt.Printf("%v\n\n", &p)

}
func TestParse1(t *testing.T) {
	b := bytes.NewBuffer([]byte(data1))
	p := parser{r: b}
	fmt.Println(p.parse())
	fmt.Printf("%v\n\n", &p)
	fmt.Println(p.parse())
	fmt.Printf("%v\n\n", &p)
	fmt.Println(p.parse())
	fmt.Printf("%v\n\n", &p)
}

func TestParseDataStream(t *testing.T) {
	b := bufio.NewScanner(bytes.NewBuffer([]byte(data)))
	result := parsedstream(b, 4)
	fmt.Println(result)
}

var data = `1 10 100 50
2 20 200 100
3 30 300 150
4 40 400 200
5 50 500 250
6 60 600 300
7 70 700 350
8 80 800 400
9 90 900 450
10 100 1000 500
11 110 1100 550
12 120 1200 600
13 130 1300 650
14 140 1400 700
15 150 1500 750
16 160 1600 800
17 170 1700 850
18 180 1800 900
19 190 1900 950
20 200 2000 1000`

var data1 = `1 10 100 50
2 20 200 100
3 30 300 150
4 40 400 200
5 50 500 250
6 60 600 300
7 70 700 350
8 80 800 400
9 90 900 450
-n
6 60 600 300
7 70 700 350
8 80 800 400
9 90 900 450
10 100 1000 500
11 110 1100 550
12 120 1200 600
13 130 1300 650
14 140 1400 700
15 150 1500 750
16 160 1600 800
17 170 1700 850
18 180 1800 900
19 190 1900 950
20 200 2000 1000`
