// basic color library with automatically expanding color palette.
// colorset.new() takes the next color in the palette and creates a new one if necessary.
package main

import "fmt"

type colorset struct {
	allcolors []color
	used      map[color]bool
}

type color uint32

func (c color) String() string {
	return fmt.Sprintf("#%06x", uint32(c))
}

const red color = 0xff0000
const green color = 0x00ff00
const blue color = 0x0000ff

func blend(c1, c2 color) color {
	const red = 0xff0000
	const green = 0x00ff00
	const blue = 0x0000ff
	r := (((c1 & red) + (c2 & red)) / 2) & red
	g := (((c1 & green) + (c2 & green)) / 2) & green
	b := (((c1 & blue) + (c2 & blue)) / 2) & blue
	return r | g | b
}
func (c *colorset) new() color {
	//initialization case
	l := len(c.allcolors)
	if l == 0 {
		c.allcolors = []color{red, green, blue}
		c.used = make(map[color]bool)
		return 0x000000 //start with black
	}
	//out of colors case, add new colors
	if l == len(c.used) {
		var cnew []color
		// for each pair of colors a_n, a_n+1 insert new blended color in the middle
		for i, _ := range c.allcolors {
			nc := blend(c.allcolors[i], c.allcolors[(i+1)%l]) //new color
			cnew = append(cnew, c.allcolors[i])               //original color
			cnew = append(cnew, nc)
		}
		c.allcolors = cnew
		return c.new()
	}
	for _, col := range c.allcolors {
		if !c.used[col] {
			c.used[col] = true
			return col
		}
	}
	fmt.Println("hit end of new() - that's not supposed to happen")
	return red
}
