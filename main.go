package main

import (
	"fmt"
	"net"
)

type Message struct {
	from string
	msg  string
}

type Server struct {
	ListenAddr string
	ln         net.Listener
	messages   chan Message
	closeChan  chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		ListenAddr: listenAddr,
		messages:   make(chan Message),
		closeChan:  make(chan struct{}),
	}
}

func (s *Server) Serve() error {
	ln, err := net.Listen("tcp", ":"+s.ListenAddr)
	if err != nil {
		return err
	}

	fmt.Println("listen: ", s.ListenAddr)
	defer ln.Close()
	s.ln = ln

	go s.AcceptConnections()
	for {
		msg := <-s.messages
		fmt.Println("from: ", msg.from, " msg: ", msg.msg)
	}

	<-s.closeChan
	close(s.closeChan)
	return nil
}

func (s *Server) AcceptConnections() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		fmt.Println("accept: ", conn.RemoteAddr())

		go s.ReadLoop(conn)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error: ", err)
			continue
		}

		msg := Message{
			from: conn.RemoteAddr().String(),
			msg:  string(buf[:n]),
		}

		s.messages <- msg
	}
}

func main() {
	server := NewServer("8080")
	server.Serve()
}
