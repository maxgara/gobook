package main

import (
	"fmt"
	"net"
	"time"
)

// func main() {
// l, err := net.Listen("tcp", "localhost:8000")
// if err != nil {
// 	fmt.Printf("listen: %v\n", err)
// 	return
// }
// for {
// 	conn, err := l.Accept()
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	go handle(conn)
// }
// }

func handle(conn net.Conn) {
	for {
		fmt.Fprintf(conn, "%s\n", time.Now().Format("01-02-06 03:04:05.00"))
		time.Sleep(time.Second)
	}
}
