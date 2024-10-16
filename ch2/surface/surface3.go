package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	width, height = 600, 500
	cells         = 80
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6 // angle of x,y axes
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/submit", ghandler) //graph handler
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// return canvas point x, y at corner of cell i,j. last return val indicates status (ok/err).
func corner(i, j int) (float64, float64, string, int) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	if math.IsNaN(z) {
		return 0, 0, "", 1
	}
	//calculate color
	rcomp := uint8(z * 255 / xyrange * 3)
	bcomp := uint8(255 - rcomp)
	str := fmt.Sprintf("%02x00%02x", rcomp, bcomp) //SVG stroke attr

	//project x,y,z onto 2-D canvas surface
	sx := width/2 + (x-y)*cos30*xyscale
	sy := width/2 + (x+y)*sin30*xyscale + -z*zscale
	return sx, sy, str, 0

}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

func f1(x, y float64) float64 {
	return x * y / 1000
}

/*
generate an SVG image and write it to out. if ref==true then the svg element tags at end+beginning will be omitted,
and only polygon tags will be returned
*/
func svggen(out io.Writer) {
	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	//<meta http-equiv="refresh" content="n">

	var maxx float64
	var maxy float64
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, astr, aerr := corner(i+1, j)
			bx, by, _, berr := corner(i, j)
			cx, cy, _, cerr := corner(i, j+1)
			dx, dy, _, derr := corner(i+1, j+1)
			//skip error polygons and continue
			if aerr|berr|cerr|derr > 0 {
				continue
			}

			fmt.Fprintf(out, "<polygon points='%.5g,%.5g %.5g,%.5g %.5g,%.5g %.5g,%.5g' stroke='#%s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, astr)
			var xs = []float64{ax, bx, cx, dx}
			var ys = []float64{ay, by, cy, dy}
			for k := range xs {
				if xs[k] > maxx {
					maxx = xs[k]
				}
				if ys[k] > maxy {
					maxy = ys[k]
				}
			}
		}

	}
	fmt.Fprintf(out, "</svg>")

	_, _ = maxx, maxy
	// fmt.Printf("max x:%f, maxy:%f", maxx, maxy)
}

// handle homepage requests
func handler(w http.ResponseWriter, r *http.Request) {
	ftext, err := ioutil.ReadFile("test.html")
	if err != nil {
		log.Fatal("handler: error reading html file")
	}
	fmt.Fprintf(w, "%s", ftext)
}

// handle eq submissions
func ghandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for f, v := range r.Form {
		fmt.Fprintf(w, "%v=%v\n", f, v)
	}
	toks := parseEq(r.Form["eq"][0])
	tokstr := SprintToks(toks)
	fmt.Println(tokstr)
	fmt.Fprintf(w, "%s", tokstr)
	const x, y = 1, 1
	val := evalEq(toks, x, y)
	fmt.Fprintf(w, "Evaluation at f(%v,%v) = %v", x, y, val)
	/******** TO DO: write eq parser********/
	// w.Header().Set("Content-Type", "image/svg+xml")
	// svggen(w)
}

const (
	VAL = iota
	VAR
	GROUP
	PLUS
	MINUS
	DIV
	MULT
	EXP
)

type Token struct {
	tt    int     //token type enum
	vname int     //only for VAR
	val   float64 //only for VAL and VAR types
	gp    []Token // subgroup token pointer, only for GROUP type
}

var legalvars = []rune{'x', 'y', 'X', 'Y', 'r', 'R'}

func parseEq(s string) []Token {
	//try s = x + y
	out := make([]Token, 0)
	var t int
	var valhit bool
	var valstr string
	var val float64
	for i, c := range s {
		if valchar(c) {
			valstr = valstr + string(c)
			valhit = true
			continue
			//when val ends, parse val str
		} else if valhit {
			var err error
			val, err = strconv.ParseFloat(valstr, 64)
			out = append(out, Token{VAL, 0, val, nil})
			if err != nil {
				log.Fatal("error: invalid str in equation" + valstr)
			}
			valhit, valstr = false, ""
		}
		if varchar(c) {
			out = append(out, Token{VAR, int(c), 0, nil})
			continue
		}
		if c == '(' {
			var nestct int
			for j := i + 1; j < len(s); j++ {
				if s[j] == ')' && nestct == 0 {
					gstr := s[i:j]
					gtokens := parseEq(gstr)
					out = append(out, Token{GROUP, 0, 0, gtokens})
					break
				} else if s[j] == ')' {
					nestct--
				}
				if s[j] == '(' {
					nestct++
				}
			}
			continue
		}
		switch c {
		case '+':
			t = PLUS
		case '-':
			t = MINUS
		case '/':
			t = DIV
		case '*':
			t = MULT
		case '^':
			t = EXP
		}
		out = append(out, Token{t, 0, 0, nil})
	}
	//clean cache for trailing VALs
	if valhit {
		var err error
		val, err = strconv.ParseFloat(valstr, 64)
		if err != nil {
			log.Fatal("error: invalid str in equation" + valstr)
		}
		out = append(out, Token{VAL, 0, val, nil})
		valhit, valstr = false, ""
	}
	return out
}
func evalEq(tstr []Token, x float64, y float64) float64 {
	var out float64
	tcount := len(tstr)
	for i := 0; i < tcount; i++ {
		t := tstr[i]
		tt := t.tt //token type
		if tt == VAR {
			vname := rune(t.vname)
			if vname > 'z' {
				vname = vname - ('A' - 'a') // getting lazy
			}
			switch vname {
			case 'x':
				out = x
			case 'y':
				out = y
			case 'r':
				out = math.Hypot(x, y)
			}
			continue
		} else if tt == VAL && i == 0 {
			out = t.val
			continue
		} else if tt == VAL {
			log.Fatal("illegal space or something - VAL hit without leading op")
			continue
		} else if tt == GROUP {
			out = evalEq(t.gp, x, y)
			continue
		}
		//handle operator tokens:
		nextval := evalEq(tstr[i+1:i+2], x, y) //can't use tstr[i+1].val because of groups
		if tt == PLUS {
			out += nextval
		} else if tt == MINUS {
			out -= nextval
		} else if tt == DIV {
			out /= nextval
		} else if tt == MULT {
			out *= nextval
		} else if tt == EXP {
			out = math.Pow(out, nextval)
		}
		//if t was an operator (no-continue case), then skip next token
		i++
	}
	return out
}

func varchar(c rune) bool {
	for _, v := range legalvars {
		if c == v {
			return true
		}
	}
	return false
}
func valchar(c rune) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	if c == '-' || c == '.' {
		return true
	}
	return false
}
func SprintToks(arr []Token) string {
	fmt.Println("TOKENS")
	var out string
	for i, v := range arr {
		out += fmt.Sprintf("arr[%d]:%v\n", i, v)
	}
	return out
}
