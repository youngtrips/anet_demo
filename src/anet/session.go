package anet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Session struct {
	id            int32
	conn          *net.TCPConn
	proto         Protocol
	wbuf          chan Message
	events        chan Event
	ctrl          chan bool
	net           string
	raddr         *net.TCPAddr
	autoReconnect bool
	reconnect     chan bool
}

const (
	SEND_BUFF_SIZE   = 1024
	CONNECT_INTERVAL = 1000 // reconnect interval
)

func newSession(id int32, conn *net.TCPConn, proto Protocol) *Session {
	sess := &Session{
		id:            id,
		conn:          conn,
		proto:         proto,
		wbuf:          make(chan Message, SEND_BUFF_SIZE),
		events:        nil,
		ctrl:          make(chan bool, 1),
		net:           "",
		raddr:         nil,
		autoReconnect: false,
		reconnect:     nil,
	}
	return sess
}

func ConnectTo(network string, addr string, proto Protocol, events chan Event, autoReconnect bool) *Session {
	session := newSession(0, nil, proto)
	session.connect(network, addr, events, autoReconnect)
	return session
}

func (self *Session) Start(events chan Event) {
	self.events = events
	go self.reader()
	go self.writer()
}

func (self *Session) ID() int32 {
	return self.id
}

func (self *Session) Close() {
	if self.autoReconnect {
		self.reconnect <- false
	}
	self.conn.Close()
}

func (self *Session) Send(api int16, payload interface{}) {
	if len(self.wbuf) < SEND_BUFF_SIZE {
		self.wbuf <- Message{api, payload}
	} else {
		log.Println("send overflow")
	}
}

func (self *Session) reader() {
	log.Printf("session[%d] start reader...", self.id)
	defer func() {
		log.Println("reader quit...")
		self.ctrl <- true
		if self.autoReconnect {
			self.reconnect <- true
		} else {
			self.events <- newEvent(EVENT_DISCONNECT, self, nil)
		}
	}()
	header := make([]byte, 2)
	for {
		if _, err := io.ReadFull(self.conn, header); err != nil {
			break
		}
		size := binary.BigEndian.Uint16(header)
		log.Printf("size=%d", size)
		data := make([]byte, size)
		if _, err := io.ReadFull(self.conn, data); err != nil {
			self.events <- newEvent(EVENT_RECV_ERROR, self, err)
			break
		}
		log.Printf("len(data)=%d", len(data))
		log.Printf("payload: %v", data)
		api, payload, err := self.proto.Decode(data)
		log.Printf("api=%v, payload=%v, err=%v", api, payload, err)
		if err != nil {
			self.events <- newEvent(EVENT_RECV_ERROR, self, err)
			break
		}
		msg := NewMessage(api, payload)
		self.events <- newEvent(EVENT_MESSAGE, self, msg)
	}
}

func encode(proto Protocol, msg Message) ([]byte, error) {
	data, err := proto.Encode(msg.Api, msg.Payload)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, uint16(len(data))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func rawSend(w *bufio.Writer, data []byte) error {
	for _, b := range data {
		fmt.Printf("%02x ", b)
	}
	fmt.Printf("\n")

	if _, err := w.Write(data); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

func (self *Session) writer() {
	log.Printf("session[%d] start writer...", self.id)
	defer func() {
		log.Println("writer quit ...")
		close(self.wbuf)
		self.conn.Close()
	}()
	w := bufio.NewWriter(self.conn)
	for {
		select {
		case msg, ok := <-self.wbuf:
			if ok {
				if raw, err := encode(self.proto, msg); err != nil {
					self.events <- newEvent(EVENT_SEND_ERROR, self, err)
					return
				} else {
					if err := rawSend(w, raw); err != nil {
						self.events <- newEvent(EVENT_SEND_ERROR, self, err)
						return
					}
				}
			} else {
				return
			}
		case <-self.ctrl:
			return
		}
	}
}

func (self *Session) supervisor() {
	defer func() {
		log.Println("supervisor quit...")
	}()
	for {
		select {
		case flag, ok := <-self.reconnect:
			if ok {
				if flag {
					log.Printf("reconnect to %s", self.raddr)
					go self.connector()
				} else {
					return
				}
			}
		}
	}
}

func (self *Session) connect(network string, addr string, events chan Event, autoReconnect bool) error {
	log.Printf("try to connect to %s %s", network, addr)
	raddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return err
	}
	self.events = events
	self.net = network
	self.raddr = raddr
	if autoReconnect {
		self.autoReconnect = autoReconnect
		self.reconnect = make(chan bool, 1)
		go self.supervisor()
	}
	go self.connector()
	return nil
}

func (self *Session) connector() {
	conn, err := net.DialTCP(self.net, nil, self.raddr)
	if err != nil {
		log.Printf("connect to %s falied: %s, id=%d", self.raddr, err, self.id)
		if self.autoReconnect {
			time.Sleep(CONNECT_INTERVAL * time.Millisecond)
			self.reconnect <- true
		} else {
			self.events <- newEvent(EVENT_CONNECT_SUCCESS, self, err)
		}
	} else {
		log.Printf("connect to %s ok...id=%d", self.raddr, self.id)
		self.conn = conn
		if !self.autoReconnect {
			self.events <- newEvent(EVENT_CONNECT_SUCCESS, self, nil)
		} else {
			self.wbuf = make(chan Message, SEND_BUFF_SIZE)
			self.Start(self.events)
		}
	}
}

func (self *Session) RemoteAddr() string {
	if self.raddr == nil {
		return ""
	}
	return self.raddr.String()
}
