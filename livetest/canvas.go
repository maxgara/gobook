// this file should define the Canvas object which contains modular html/svg/other output objects as
// well as input objects.
package main

import (
	"fmt"
	"strings"
)

// Fmap dest specifiers:
// $APPEND
// $PREPEND
type Fmap struct {
	src  string              //id of caller item
	f    func(string) string //function mapped to that path
	dest string              //id of item to re-write, or a flag

}

// A Canvas object displays golang objects in a browser-viewable format, and maps input objects to functions
// modifying these objects, or the Canvas itself. Behaviour is controlled by use of the Canvas.conf string.

type Canvas struct {
	id      string          //name provided by Canvas caller
	items   []*fmt.Stringer //go objects to be displayed on canvas
	itemIds []string        //ids for go objects displayed by Canvas
	ops     []Fmap          //supported input function calls for frontend
	conf    canvasConfig    //all configuration settings
}

func NewCanvas(id string) *Canvas {
	c := Canvas{id: id}
	return &c
}

// adds new item to c
func (c *Canvas) NewItem(id string, obj fmt.Stringer) {
	c.items = append(c.items, &obj)
	c.itemIds = append(c.itemIds, id)
}

// HTML String
func (c *Canvas) String() string {
	var itemstrs []string
	for idx := range c.items {
		item := *(c.items[idx])
		itemid := c.itemIds[idx]
		div := webStringify(item, itemid, c)
		itemstrs = append(itemstrs, div)
	}
	return canvasStrsWrapper(itemstrs)
}

// get a "webstring" from item, either using item.String directly or with adjusted formatting.
// also wrap the string in a div for tracking.
func webStringify(item fmt.Stringer, itemid string, c *Canvas) string {
	var itemstr string
	// //check if obj does it's own web formatting
	if c.conf.Get("itemconfs."+itemid+".WebStringer") != nil {
		itemstr = item.String()
	} else {
		itemstr = stringToHTML(item.String())
	}
	div := itemDiv(itemid, itemstr)
	return div
}

// put item in div
func itemDiv(id string, s string) string {
	return fmt.Sprintf("<div id=\"%s\">%s</div>", id, s)
}

// wrap items in divs into html template
func canvasStrsWrapper(itemstrs []string) string {
	fstr := `<html><head> <script src="https://unpkg.com/htmx.org@2.0.3"></script></head><body>%s</body></html>`
	var s string
	for _, item := range itemstrs {
		s += item
	}
	return fmt.Sprintf(fstr, s)
}

// replace special characters in string with HTML equivalents
func stringToHTML(s string) string {
	s = strings.ReplaceAll(s, " ", "&nbsp")
	s = strings.ReplaceAll(s, "\n", "</br>")
	return s
}

func (t *textAreaStruct) String() string {
	return t.str
}

// used for RegisterInput function
type textAreaStruct struct {
	id  string
	str string
}

func (c *Canvas) NewInputTextArea(id string, f func(string) string, targetid string) {
	fstr := `<label for="%s">%[1]s</label>
		<textarea hx-post="/apply/" hx-trigger="keyup delay:500ms changed" type="text" id="%[1]s" name="input"
			hx-target="#%s"></textarea>`
	s := fmt.Sprintf(fstr, id, targetid)
	item := textAreaStruct{id: id, str: s}
	c.NewItem(id, &item)
	c.conf.Set("itemconfs."+id+".WebStringer", "t")            //tell canvas not to reformat
	c.ops = append(c.ops, Fmap{src: id, f: f, dest: targetid}) //register handler func
}

// apply input from user to relevant item, send back a string representation of that item's new state
func (c *Canvas) Apply(id string, arg string) string {
	for _, m := range c.ops {
		if m.src == id {
			m.f(arg) //call function
			target := *ItemForId(c, m.dest)
			return webStringify(target, m.dest, c)
		}
	}
	panic("no handler for Apply call")
}

func ItemForId(c *Canvas, id string) *fmt.Stringer {
	for i, v := range c.itemIds {
		if v == id {
			return c.items[i]
		}
	}
	return nil //ID not listed, error condition
}
