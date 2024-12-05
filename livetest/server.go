package main

// import (
// 	"fmt"
// 	"html/template"
// 	"io"
// 	"log"
// 	"net/http"
// )

// var doctempl *template.Template
// var statetempl *template.Template
// var textareatempl *template.Template

// const docTemplFile = "/Users/maxgara/Desktop/go-code/gobook/workspace/livetest/basic.tmpl"
// const stateTemplFile = "/Users/maxgara/Desktop/go-code/gobook/workspace/livetest/state.tmpl"
// const textareaTemplFile = "/Users/maxgara/Desktop/go-code/gobook/workspace/livetest/textarea.tmpl"

// // web module interface. Implementation should not be a reference type, as the value of the wm
// // used to call Live() is used as an initialization constant. WString should return an HTML string that can
// // be displayed in a <div> element within a webpage. Should not modify wm.
// // WInput should process a line of text input in a way that is transparent to the user.

// type WebModule interface {
// 	WString() string
// 	WInput(string)
// 	WInit()
// }

// type WebModuleD interface {
// 	WebModule
// 	FirstChild() bool
// 	NextSibling() bool
// }

// type tempstruct struct {
// 	Webstring string
// }

// // input state
// type state struct {
// 	calls map[string]string //map of [input_id]input
// 	order []string          //[]input_id
// }

// type uicomp struct {
// 	Id string //html id
// }

// // get system state for input state
// func runstate(wm WebModule, st state) {
// 	for _, id := range st.order {
// 		wm.WInput(st.calls[id])
// 	}
// }

// // process input into calls and add to struct
// func statechange(stt *state, call string, caller string) {
// 	fmt.Printf("state change: str=%v\n", call)
// 	//add new caller if missing
// 	if _, ok := stt.calls[caller]; !ok {
// 		stt.order = append(stt.order, caller)
// 	}
// 	stt.calls[caller] = call //update call in state
// }

// // serve object's WString within a basic HTML template, server address is localhost:8000.
// func liveprint(wm WebModule) {
// 	//initialize
// 	wm.WInit()
// 	st := state{calls: make(map[string]string)}
// 	var ui []uicomp
// 	//serve document
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("doc root handler called.\n")
// 		doctempl.Execute(w, struct{}{}) //write html to resp
// 		fmt.Printf("sent doc, handler done.\n")
// 	}
// 	//serve current state to user
// 	readHandler := func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("read handler called. state is %v\n", st)
// 		wm.WInit()
// 		runstate(wm, st) //get system state for input state
// 		fmt.Printf("after runstate, webmodule obj= %v\n", wm)
// 		obj := tempstruct{Webstring: wm.WString()}
// 		fmt.Printf("getting WString from webmodule:%v\n", wm.WString())
// 		statetempl.Execute(w, obj) //write html to resp
// 		fmt.Printf("sent resp. handler done.\n")
// 	}
// 	//update state to match user
// 	inputHandler := func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("input handler called. current (unsynced) server state :%v\n", st)
// 		r.ParseForm()
// 		callstr := r.Form.Get("input")
// 		trigger := r.Header.Get("HX-Trigger")
// 		statechange(&st, callstr, trigger)
// 		fmt.Printf("input handled. state:%v\n", st)
// 		fmt.Printf("input handler calling read handler.\n")
// 		readHandler(w, r) // send back updated state to user
// 	}
// 	uiUpdateHandler := func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("ui update handler called: current ui: %v\n", ui)
// 		r.ParseForm()
// 		_ = r.Form.Get("HX-Trigger")
// 		newcomp := uicomp{"input" + fmt.Sprint(len(ui))}
// 		ui = append(ui, newcomp)
// 		textareatempl.Execute(w, newcomp)
// 	}

// 	doctempl = template.Must(template.ParseFiles(docTemplFile))
// 	statetempl = template.Must(template.ParseFiles(stateTemplFile))
// 	textareatempl = template.Must(template.ParseFiles(textareaTemplFile))
// 	http.HandleFunc("/input", inputHandler)
// 	http.HandleFunc("/", handler)
// 	http.HandleFunc("/output", readHandler)
// 	http.HandleFunc("/uiupdate", uiUpdateHandler)
// 	log.Fatal(http.ListenAndServe("localhost:8000", nil))
// }
// func printr(wm WebModuleD, w io.Writer) {
// 	fmt.Printf("getting WString from webmodule:%v\n", wm.WString())
// 	obj := tempstruct{Webstring: wm.WString()}
// 	fmt.Printf("writing current wm")
// 	statetempl.Execute(w, obj) //write html to resp
// 	fmt.Printf("moving on to descs\n")
// 	if wm.FirstChild() {
// 		fmt.Println("got child")
// 		printr(wm, w)
// 	}
// 	if wm.NextSibling() {
// 		fmt.Println("got sib")
// 		printr(wm, w)
// 	}
// }

// func liveprintd(wm WebModuleD) {
// 	//initialize
// 	wm.WInit()
// 	st0 := state{calls: make(map[string]string)}
// 	st := st0
// 	var ui []uicomp
// 	//serve document
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		st = st0
// 		doctempl = template.Must(template.ParseFiles(docTemplFile))
// 		statetempl = template.Must(template.ParseFiles(stateTemplFile))
// 		textareatempl = template.Must(template.ParseFiles(textareaTemplFile))
// 		fmt.Printf("doc root handler called.\n")
// 		doctempl.Execute(w, struct{}{}) //write html to resp
// 		fmt.Printf("sent doc, handler done.\n")
// 	}
// 	//serve current state to user
// 	readHandler := func(w http.ResponseWriter, r *http.Request) {
// 		doctempl = template.Must(template.ParseFiles(docTemplFile))
// 		statetempl = template.Must(template.ParseFiles(stateTemplFile))
// 		textareatempl = template.Must(template.ParseFiles(textareaTemplFile))
// 		fmt.Printf("read handler called. state is %v\n", st)
// 		wm.WInit()
// 		runstate(wm, st) //get system state for input state
// 		printr(wm, w)
// 		fmt.Printf("sent resp. handler done.\n")
// 	}
// 	//update state to match user
// 	inputHandler := func(w http.ResponseWriter, r *http.Request) {
// 		doctempl = template.Must(template.ParseFiles(docTemplFile))
// 		statetempl = template.Must(template.ParseFiles(stateTemplFile))
// 		textareatempl = template.Must(template.ParseFiles(textareaTemplFile))
// 		fmt.Printf("input handler called. current (unsynced) server state :%v\n", st)
// 		r.ParseForm()
// 		callstr := r.Form.Get("input")
// 		trigger := r.Header.Get("HX-Trigger")
// 		statechange(&st, callstr, trigger)
// 		fmt.Printf("input handled. state:%v\n", st)
// 		fmt.Printf("input handler calling read handler.\n")
// 		readHandler(w, r) // send back updated state to user
// 	}
// 	uiUpdateHandler := func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("ui update handler called: current ui: %v\n", ui)
// 		r.ParseForm()
// 		_ = r.Form.Get("HX-Trigger")
// 		newcomp := uicomp{"input" + fmt.Sprint(len(ui))}
// 		ui = append(ui, newcomp)
// 		textareatempl.Execute(w, newcomp)
// 	}

// 	doctempl = template.Must(template.ParseFiles(docTemplFile))
// 	statetempl = template.Must(template.ParseFiles(stateTemplFile))
// 	textareatempl = template.Must(template.ParseFiles(textareaTemplFile))
// 	http.HandleFunc("/input", inputHandler)
// 	http.HandleFunc("/", handler)
// 	http.HandleFunc("/output", readHandler)
// 	http.HandleFunc("/uiupdate", uiUpdateHandler)
// 	log.Fatal(http.ListenAndServe("localhost:8000", nil))
// }

// func LivePrint(target any) {
// 	if wm, ok := target.(WebModuleD); ok {
// 		fmt.Printf("going into recursive WM mode...\n")
// 		liveprintd(wm)
// 		return
// 	}
// 	fmt.Printf("no recursive WM mode.\n")
// 	if wm, ok := target.(WebModule); ok {
// 		liveprint(wm)
// 		return
// 	}
// 	log.Fatal("bad webmodule. Type assertion failed.")
// }
