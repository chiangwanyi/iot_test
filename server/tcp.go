package server

import (
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type TcpServer struct {
	listener net.Listener
	conns    map[net.Conn]struct{}
	connsMu  sync.Mutex
	wg       sync.WaitGroup
}

func NewTcpServer(addr string) (*TcpServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	log.Printf("Server listening on %s", addr)
	return &TcpServer{
		listener: listener,
		conns:    make(map[net.Conn]struct{}),
	}, nil
}

func (s *TcpServer) StartTcpServer() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("Server is shutting down...")
				return
			}
			log.Printf("Accept error: %v", err)
			continue
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		s.wg.Add(1)
		s.connsMu.Lock()
		s.conns[conn] = struct{}{}
		s.connsMu.Unlock()

		go s.handleConn(conn)
	}
}

func (s *TcpServer) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		s.connsMu.Lock()
		delete(s.conns, conn)
		s.connsMu.Unlock()
		s.wg.Done()
		log.Printf("Connection from %s closed", conn.RemoteAddr())
	}()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("Read error from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}
		if n > 0 {
			data := buf[:n]
			hexString := hex.EncodeToString(data)
			log.Printf("Received %d bytes from %s: %s", n, conn.RemoteAddr(), hexString)
		}
	}
}

func (s *TcpServer) StopTcpServer() {
	log.Println("Initiating graceful shutdown...")
	s.listener.Close()

	s.connsMu.Lock()
	for conn := range s.conns {
		conn.Close()
	}
	s.connsMu.Unlock()

	s.wg.Wait()
	log.Println("Graceful shutdown complete")
}
