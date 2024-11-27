package live

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var templ *template.Template

// web string interface. Any WebModule implementation wm should be a reference type.
// WString should return an HTML string that can be displayed in a <div> element within a webpage. Should not modify wm.
// WInput should process a line of text input in a way that is transparent to the web user.
// WInit should initialize a wm implementation state in a way that is invarient to previous state. (maybe start with make(<wm>))

type WebModule interface {
	WString() string
	WInput(string)
	WInit()
}

type tempstruct struct {
	Webstring string
}
type state struct {
	calls []string //input from current state of webpage
}

// create a WebModule state equivalent corresponding to st
func runstate(wm WebModule, st state) {
	wm.WInit()
	for _, str := range st.calls {
		wm.WInput(str)
	}
}

// process input into calls and add to struct
func addcall(stt *state, str string) {
	fmt.Printf("addcall: str=%v\n", str)
	str, _ = url.QueryUnescape(str)
	str = strings.TrimPrefix(str, "input=")
	stt.calls = append(stt.calls, str)
}

// serve object's WString within a basic HTML template, server address is localhost:8000.
func liveprint(target WebModule) {
	var st state
	//serve current state to user
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("handler called. state: %v\n", st)
		runstate(target, st)
		fmt.Printf("after runstate, webmodule obj= %v\n", target)
		obj := tempstruct{Webstring: target.WString()}
		fmt.Printf("getting WString from webmodule:%v\n", target.WString())
		templ.Execute(w, obj) //write html to resp
		fmt.Printf("sent resp. handler done.\n")
	}
	//update state to match user
	inputHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("input handler called. current (unsynced) server state :%v\n", st)
		bytes, _ := io.ReadAll(r.Body)
		callstr := string(bytes)
		st = state{} // reset state. Eventually may not do this if input HTTP requests do not include entire user state.
		addcall(&st, callstr)
		fmt.Printf("input handled. state:%v\n", st)
		fmt.Printf("input handler calling handler.\n")
		handler(w, r) // send back updated state to user
	}
	var err error
	templ, err = template.ParseFiles("/Users/maxgara/Desktop/go-code/gobook/workspace/livetest/basic.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/input", inputHandler)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
func LivePrint(target any) {
	if wm, ok := target.(WebModule); !ok {
		log.Fatal("bad webmodule. Type assertion failed.")
	} else {
		liveprint(wm)
	}
}
