#!/usr/bin/env python
#!coding=utf8

import struct
import types
from protocol import user_pb2
from protocol import common_pb2
from protocol import opcode_pb2
from protocol import error_pb2

"""
x	pad byte	no value	 	 
c	char	string of length 1	1	 
b	signed char	integer	1	(3)
B	unsigned char	integer	1	(3)
?	_Bool	bool	1	(1)
h	short	integer	2	(3)
H	unsigned short	integer	2	(3)
i	int	integer	4	(3)
I	unsigned int	integer	4	(3)
l	long	integer	4	(3)
L	unsigned long	integer	4	(3)
q	long long	integer	8	(2), (3)
Q	unsigned long long	integer	8	(2), (3)
f	float	float	4	(4)
d	double	float	8	(4)
s	char[]	string	 	 
p	char[]	string	 	 
P	void *	integer	 	(5), (3)
"""
def pack(opcode, msg):
    data = msg.SerializeToString()
    return struct.pack('>HH%ds' % len(data), len(data) + 2, opcode, data)

def unpack(payload):
    (opcode,) = struct.unpack('>H', payload[0:2])
    msg = None
    if opcode ==  opcode_pb2.COMMON_ACK:
        msg = common_pb2.CommonAck()
        msg.ParseFromString(payload[2:])
    elif opcode ==  opcode_pb2.LOGIN_ACK:
        msg = user_pb2.LoginAck()
        msg.ParseFromString(payload[2:])
    elif opcode == opcode_pb2.ONLINE_USERS_ACK:
        msg = user_pb2.OnlineUserList()
        msg.ParseFromString(payload[2:])
    elif opcode == opcode_pb2.EXEC_CMD_NTF:
        msg = user_pb2.ExecCmdInfo()
        msg.ParseFromString(payload[2:])
    return opcode, msg
