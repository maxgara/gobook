package main

import (
	"maxgara-code.com/workspace/garachat/rpcserver"
)

func main() {
	rpcserver.Run()
	// //initialize
	// svr := new(rpcserver.ChatServer)
	// rpc.Register(svr)
	// rpc.HandleHTTP()
	// l, err := net.Listen("tcp", ":1234")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// //wait for shutdown signal or server error
	// ch := make(chan string)
	// go func() {
	// 	ch <- http.Serve(l, nil).Error()
	// }()
	// go rpcserver.CmdRead(ch)
	// <-ch
}
