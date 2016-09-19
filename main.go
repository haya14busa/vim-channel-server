package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", os.Getenv("VIM_LISTEN_ADDRESS"), "server address")
}

func main() {
	flag.Parse()

	r, w := os.Stdin, os.Stdout

	if addr == "" {
		addr = ":0"
	}

	n := resolveNet(addr)
	if n == "unix" {
		os.Remove(addr)
	}
	ln, err := net.Listen(n, addr)
	if err != nil {
		log.Println(err)
		time.Sleep(1 * time.Second)
		log.Fatal(err)
	}
	ex(os.Stdout, fmt.Sprintf("let g:vim_channel_server#addr = '%s'", ln.Addr().String()))

	conns := map[net.Addr]net.Conn{}
	var connsMu sync.RWMutex

	// readline from stdin and send message to each connection.
	go func(s *bufio.Scanner) {
		for s.Scan() {
			connsMu.Lock()
			for addr, conn := range conns {
				go func(addr net.Addr, conn net.Conn) {
					if _, err := conn.Write(s.Bytes()); err != nil {
						delete(conns, addr)
						log.Printf("disconnected %v", addr)
					}
					conn.Write([]byte("\n"))
				}(addr, conn)
			}
			connsMu.Unlock()
		}
	}(bufio.NewScanner(r))

	// Loop for accepting new connection.
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		connsMu.Lock()
		conns[conn.RemoteAddr()] = conn
		connsMu.Unlock()
		go func(conn net.Conn) {
			if _, err := io.Copy(w, conn); err != nil {
				log.Printf("fail to send message to Vim from %v: err", conn.RemoteAddr(), err)
			}
		}(conn)
	}

	log.Println("bye;)")
}

func ex(w io.Writer, cmd string) error {
	return json.NewEncoder(w).Encode([]interface{}{"ex", cmd})
}

// resolveNet returns "tcp" or "unix" based on addr.
func resolveNet(addr string) string {
	if _, err := net.ResolveTCPAddr("tcp", addr); err == nil {
		return "tcp"
	}
	if _, err := net.ResolveUnixAddr("unix", addr); err == nil {
		return "unix"
	}
	return ""
}
