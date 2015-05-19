package anet

import (
	"log"
	"net"
	"time"
)

const (
	ID_POOL_SIZE = 32
)

type Server struct {
	net    string
	addr   string
	ln     *net.TCPListener
	eq     chan Event
	idpool chan int32
	proto  Protocol
}

func NewServer(net string, addr string, proto Protocol, eq chan Event) *Server {
	srv := Server{
		net:    net,
		addr:   addr,
		ln:     nil,
		eq:     eq,
		idpool: make(chan int32, ID_POOL_SIZE),
		proto:  proto,
	}
	return &srv
}

func (s *Server) ListenAndServe() error {
	tcpAddr, err := net.ResolveTCPAddr(s.net, s.addr)
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP(s.net, tcpAddr)
	if err != nil {
		return err
	}
	go func() {
		id := int32(1)
		for {
			s.idpool <- id
			id++
		}
	}()
	go func() {
		defer ln.Close()
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, e := ln.AcceptTCP()
			if e != nil {
				if ne, ok := e.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					log.Printf("http: Accept error: %v; retrying in %v", e, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				break
			}
			tempDelay = 0
			id := s.nextID()
			session := newSession(id, conn, s.proto)
			s.eq <- newEvent(EVENT_ACCEPT, session, nil)
		}
	}()
	s.ln = ln
	return nil
}

func (s *Server) nextID() int32 {
	return <-s.idpool
}

func (s *Server) Close() {
	s.ln.Close()
}
