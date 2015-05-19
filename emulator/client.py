#!/usr/bin/env python
#!coding=utf8

import sys
import hashlib
import struct
import session
from protocol import user_pb2
from protocol import common_pb2
from protocol import opcode_pb2
from protocol import error_pb2

import packet

class Client(object):
    def __init__(self):
        object.__init__(self)
        self.session = session.Session()
        self.uid = None
        self.username = ""

    def open(self, ip, port):
        return self.session.open(ip, port)

    def start(self):
        self.session.start()

    def stop(self):
        self.session.stop()

    def join(self):
        self.session.join()

    def login(self, username, password):
        self.username = username
        #self.session.send(username + ':' + password)
        req = user_pb2.LoginReq()
        req.username = username
        req.password = password
        md5 = hashlib.md5()
        md5.update(req.password)
        req.password = md5.hexdigest()
        self.session.send(packet.pack(opcode_pb2.LOGIN_REQ, req))

    def logout(self):
        pass

    def register(self, username, password):
        self.username = username

    def onlines(self):
        req = common_pb2.NullMessage()
        self.session.send(packet.pack(opcode_pb2.ONLINE_USERS_REQ, req))

    def exec_cmd(self, userId, cmd):
        req = user_pb2.ExecCmdReq()
        req.cmd = cmd
        req.target_uid = userId
        self.session.send(packet.pack(opcode_pb2.EXEC_CMD_REQ, req))
