package main

// open RPC conn with rpcserver. Send messages and poll for responses every 2 sec.

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"time"

	"maxgara-code.com/workspace/garachat/rpcserver"
)

func main() {
	var errors int
	for {
		err := connect()
		fmt.Println(err.Error() + "\n")
		errors++
	}
}

// initialize connection
func connect() error {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	var msgs []rpcserver.Message
	var msgReadIdx int //message index
	var name string    //username
	fmt.Println("please enter your username")
	fmt.Scanln(&name)
	fmt.Println("enter your password or press enter now to create a new user")
	var pw uint64
	fmt.Scanln(&pw)
	if pw == 0 {
		err := client.Call("ChatServer.NewUser", rpcserver.Args{Usr: name}, &pw)
		if err != nil {
			return fmt.Errorf("rpc Create call for user %v: %v", name, err)
		}
		fmt.Printf("pw created:'%d'\n", pw)
	}
	fmt.Println("---CHAT---")
	client.Call("ChatServer.ReadLast", rpcserver.Args{T: pw, Usr: name, N: 1000}, &msgs) //should never return error
	if err != nil {
		return fmt.Errorf("readlast err: %v", err)
	}
	for _, msg := range msgs {
		fmt.Printf("%v\n", msg)
	}
	//poll for updates
	go func() {
		for {
			time.Sleep(1000000 * 2) //two second polling delay
			client.Call("ChatServer.ReadFrom", rpcserver.Args{Usr: name, T: pw, Idx: msgReadIdx}, &msgs)
			//print new messages
			for _, msg := range msgs {
				fmt.Printf("%v\n", msg)
			}
			msgReadIdx += len(msgs)
		}
	}()
	//send messages to server from stdin
	go func() {
		const BSIZE = 1000
		var bytes = make([]byte, BSIZE)
		var i int
		for {
			//read stdin to bytes until newline or EOF or length limit
			b := bytes[i:i+1]
			n, err := os.Stdin.Read(b)
			if err != io.EOF && b != '\n'{
				if err != nil{
					panic()
				}
				if i+1 < BSIZE{
					i++ 
					continue
				}
				fmt.Printf("\nauto-sending message: buffer full\n")
			}
			//send message to chat and start over
			m := string(bytes[:i+n]) // include last byte read, usually newline
			var resp string
			var args = rpcserver.Args{name, pw, msgReadIdx, 0, m}
			err = client.Call("ChatServer.Submit", &args, &resp)
			if err != nil {
				fmt.Printf("submit %v: %v\n", args, err)
			}
			i = 0
		}
	}()
	// wait forever
	c := make(chan int)
	c <- 0
	return nil
}
