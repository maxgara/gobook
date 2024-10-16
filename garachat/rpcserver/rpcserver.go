package rpcserver

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	"sync"
)

// allow chat client RPC to create users, send chat messsages, and read the chat log
type ChatServer struct {
	creds map[string]uint64 //user creds
	msgs  []Message         // message log
	mu    sync.Mutex
}
type Message struct {
	Usr string
	Msg string
}
type Args struct {
	Usr   string //username
	Token uint64 //pseudo auth token
	Idx   int    //read index in chatlog
	N     int    // number of messages to retreive when calling ReadLast
	Msg   string //used for sending a new message
}

func (m Message) String() string {
	return fmt.Sprintf("%s: %s", m.Usr, m.Msg)
}

// authenticate user
func (s *ChatServer) auth(args Args) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	key, ok := s.creds[args.Usr]
	if !ok {
		fmt.Printf("INFO: failed auth: nonexistent user %s\n", args.Usr)
		return false
	}
	if key != args.Token {
		fmt.Printf("INFO: failed auth: bad password for user %s. Creds:%v\n", args.Usr, s.creds)
		return false
	}
	return true
}

// rpc func - create user
func (s *ChatServer) NewUser(args Args, resp *uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.creds[args.Usr]; ok {
		return fmt.Errorf("user %v already exists", args.Usr)
	}
	if strings.ContainsAny(args.Usr, " \t\n\r") {
		return fmt.Errorf("no whitespace or newlines allowed:'%v'", args.Usr)
	}
	if args.Usr == "" {
		return fmt.Errorf("username must not be empty")
	}
	t := rand.Uint64()
	fmt.Printf("token %d generated for user %s\n", t, args.Usr)
	s.creds[args.Usr] = t
	*resp = t
	return nil
}

// rpc func - submit message
func (s *ChatServer) Submit(args Args, resp *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.auth(args) {
		return fmt.Errorf("bad username or password. Please try again")
	}
	m := Message{args.Usr, args.Msg}
	s.msgs = append(s.msgs, m)
	ms := m.Msg
	//remove newlines
	if ms[len(ms)-1] == '\n' {
		ms = ms[:len(ms)-2]
	}
	fmt.Printf("INFO:CHAT \"%v\"", ms)
	return nil
}

// rpc func - get chatlog updates from server
func (s *ChatServer) ReadFrom(args Args, resp *[]Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.auth(args) {
		return fmt.Errorf("bad username or password. Please try again")
	}
	if args.Idx >= len(s.msgs) {
		return nil // no new messages
	}
	*resp = s.msgs[args.Idx:]
	return nil
}
func (s *ChatServer) ReadLast(args Args, resp *[]Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.auth(args) {
		return fmt.Errorf("bad username or password. Please try again")
	}
	off := len(s.msgs) - args.N
	if off < 0 {
		off = 0
	}
	*resp = s.msgs[off:]
	return nil
}

// listen for 'q' key to quit
func CmdRead(ch chan string) {
	buff := make([]byte, 1)
	for {
		os.Stdin.Read(buff)
		if buff[0] == 'q' {
			break
		}
	}
	ch <- "server shutdown from cmdline"
}
func InitServer() *ChatServer {
	//initialize
	svr := new(ChatServer)
	svr.creds = make(map[string]uint64)
	svr.msgs = make([]Message, 0)
	return svr
}

// start a chat server instance, register, listen and serve RPC calls on port :1234
func Run() {
	svr := InitServer()
	rpc.Register(svr)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	//wait for shutdown signal or server error
	ch := make(chan string)
	go func() {
		ch <- http.Serve(l, nil).Error()
	}()
	go CmdRead(ch)
	<-ch
}
