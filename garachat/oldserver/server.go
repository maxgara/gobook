package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

type message struct {
	Usr string //message sender username
	Txt string //message content
}
type messagelog struct {
	Msgs []message
}

func (m message) String() string {
	return fmt.Sprintf("[%s]:%s\n", m.Usr, m.Txt)
}

func (l messagelog) String() string {
	s := ""
	for _, msg := range l.Msgs {
		s += msg.String()
	}
	return s
}

// add message to log
func (l *messagelog) Add(m message) {
	l.Msgs = append(l.Msgs, m)
}

// chat server
type msgserver struct {
	msglg *messagelog
	templ *template.Template
	conns []*websocket.Conn
}

// initialize, load chat log, compile HTML template
func (s *msgserver) initSvr() {
	var msglg messagelog
	msglg.load() // load chat history from file
	s.conns = []*websocket.Conn{}
	s.templ = template.Must(template.New("homepage").Parse(homepageStr)) //compile template for html
}

func main() {
	svr := msgserver{} //create message server
	svr.initSvr()
	defer svr.msglg.save() //save chat history at server shutdown
	http.HandleFunc("/", svr.mainHndlr)
	// var mh = msgHandler{&msglg, conns}
	http.Handle("/soc", websocket.Handler(svr.msgHndlr))
	ch := make(chan string)
	go asyncServe("localhost:8000", ch)
	go cmdRead(ch)
	fmt.Println(<-ch) // either server error or user requested shutdown from cmdline
}

// serves pages, sends any server error to ch
func asyncServe(url string, ch chan string) {
	err := http.ListenAndServe(url, nil) // returns on error
	ch <- err.Error()
}

// listen for 'q' key to quit
func cmdRead(ch chan string) {
	buff := make([]byte, 256)
	for {
		os.Stdin.Read(buff)
		if buff[0] == 'q' {
			break
		}
	}
	ch <- "server shutdown from cmdline"
}

// serve html
func (s *msgserver) mainHndlr(w http.ResponseWriter, r *http.Request) {
	s.templ.Execute(w, s.msglg) // insert messagelog data into html template, write to client conn
}

// handle sending and recieving messages from client (implements websocket.Handler)
func (s *msgserver) msgHndlr(ws *websocket.Conn) {
	s.conns = append(s.conns, ws)
	fmt.Printf("subscribed %p to updates\n", ws)
	for {
		err := readMsg(s.msglg, ws)
		if err != nil {
			ws.Close()
			s.removeConn(ws)
			fmt.Printf("removing %p from conns\n", ws)
			return //drop websocket conn
		}
		s.sendUpdates()
	}
}

func (s *msgserver) removeConn(c *websocket.Conn) {
	// newlen := len(s.conns) - 1
	var hit bool
	for i := range s.conns {
		if s.conns[i] == c {
			s.conns = append(s.conns[:i], s.conns[i:]...)
		}
	}
	if !hit {
		fmt.Fprintf(os.Stderr, "removeConn warning: no match for conn %p in conns\n", c)
		return
	}
}

// deliver new chatlog to all open conns
func (s msgserver) sendUpdates() {
	for _, c := range s.conns {
		// fmt.Println("sending updates to:")
		// fmt.Println(h.conns)
		// fmt.Println(h.msglg.String())
		c.Write([]byte(s.msglg.String()))
	}
}

// add message from client request to message log
func readMsg(msglg *messagelog, ws *websocket.Conn) error {
	var buff = make([]byte, 512)
	var bytes = make([]byte, 0, 512)
	for {
		n, err := ws.Read(buff)
		bytes = append(bytes, buff[:n]...)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	if len(bytes) == 0 {
		return nil
	}
	msglg.Add(message{"x", string(bytes)})
	return nil
}

func (msglg *messagelog) save() {
	f, err := os.OpenFile("chatlog.txt", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	m, err := json.Marshal(msglg)
	if err != nil {
		log.Fatal(fmt.Errorf("Marshal:%v", err))
	}
	_, err = f.Write(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "savefile write:%v", err)
		return
	}
	fmt.Printf("log file saved successfully\n")
}
func (msglg *messagelog) load() {
	data, err := os.ReadFile("chatlog.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ReadFile chatlog.txt: %v", err)
		return
	}
	if !json.Valid(data) {
		fmt.Printf("log file invalid: new log created\n")
		return //use empty messagelog if load fails
	}
	err = json.Unmarshal(data, &msglg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "savefile load:%v", err)
		return
	}
	fmt.Printf("log file loaded successfully\n")
	fmt.Println(msglg)
}
