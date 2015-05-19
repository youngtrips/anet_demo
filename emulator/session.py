#!/usr/bin/env python
#coding=utf8

import threading
import socket
import Queue
import select
import sys
import struct
import packet

from protocol import user_pb2
from protocol import common_pb2
from protocol import opcode_pb2
from protocol import error_pb2

class Session(threading.Thread):
    def __init__(self):
        threading.Thread.__init__(self)
        self.s = None
        self.q = Queue.Queue(1024)
        self.buff = ""
        self.ev = threading.Event()

    def open(self, ip, port):
        self.ip = ip
        self.port = port
        self.s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        code = self.s.connect_ex((ip, port))
        if code != 0:
            self.s.close()
        return code

    def send(self, data):
        self.q.put(data, block=False)

    def run(self):
        self.loop()

    def stop(self):
        self.ev.set()

    def loop(self):
        self.rbuf = ""
        while True:
            if self.ev.isSet():
                break
            rset, wset, _ = select.select([self.s], [self.s], [], 0.1)
            if len(rset) > 0:
                data = self.s.recv(8192)
                if len(data) == 0:
                    break
                self.rbuf += data
                opcode, msg = self.parse()
                if opcode == opcode_pb2.LOGIN_ACK:
                    self.on_login(msg)
                elif opcode == opcode_pb2.COMMON_ACK:
                    self.on_commonAck(msg)
                else:
                    print (opcode, msg)

            if len(wset) > 0:
                while not self.q.empty():
                    data = self.q.get(block=False)
                    self.buff += data
                if len(self.buff) > 0:
                    size = self.s.send(self.buff)
                    if size > 0:
                        self.buff = self.buff[size:]
        print 'disconnect to server %s:%d' % (self.ip, self.port)
        self.s.close()

    def on_login(self, msg):
        print msg

    def on_commonAck(self, msg):
        print msg

    def parse(self):
        if len(self.rbuf) < 2:
            return 0, None
        #print len(self.rbuf)
        (size,) = struct.unpack('>H', self.rbuf[0:2])
        if len(self.rbuf) < size - 2:
            return 0, None
        #print 'size=%d' % (size)
        opcode, msg = packet.unpack(self.rbuf[2:size+2])
        self.rbuf = self.rbuf[size+2:]
        return opcode, msg
