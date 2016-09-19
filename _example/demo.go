// $ VIM_LISTEN_ADDRESS=/tmp/vim.sock vim -Nu NORC --cmd 'set runtimepath+=.'
// $ go run ./_example/demo.go
// ["expr", "1+1", -1]
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("unix", "/tmp/vim.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		s := bufio.NewScanner(conn)
		for s.Scan() {
			fmt.Printf("receive: %v\n", s.Text())
		}
	}()

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		conn.Write(s.Bytes())
	}
}
