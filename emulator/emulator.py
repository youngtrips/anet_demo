#!/usr/bin/env python
#coding=utf8

import argparse
import cmd
import sys
import os

import readline
import rlcompleter

import client

import platform

if platform.system() != "Windows":
    if 'libedit' in readline.__doc__:
        readline.parse_and_bind("bind ^I rl_complete")
    else:
        readline.parse_and_bind("tab: complete")

class Args(object):
    pass


class CLI(cmd.Cmd):
    def __init__(self, clt):
        cmd.Cmd.__init__(self)
        self.prompt=">>>"
        self.clt = clt

    def emptyline(self):
        return

    def parse(self, data):
        if data == None:
            return []
        return  data.split(' ')

    def do_login(self, data):
        """login username password"""
        args = self.parse(data)
        if len(args) < 2:
            print """login username password"""
            return
        self.clt.login(args[0], args[1])

    def do_register(self, data):
        """register newuser"""
        args = self.parse(data)
        if len(args) < 2:
            print """login username password"""
            return
        self.clt.register(args[0], args[1])

    def do_onlines(self, _):
        self.clt.onlines()

    def do_exec_cmd(self, data):
        args = self.parse(data)
        if len(args) < 2:
            print """exec_cmd userId cmd"""
            return
        self.clt.exec_cmd(int(args[0]), args[1])


    def do_quit(self, _):
        """logout and quit emulator"""
        self.clt.quit(data)
        return True

    def do_q(self, _):
        """logout and quit emulator"""
        self.clt.stop()
        return True

    def do_h(self, data):
        """help [command]"""
        self.do_help(data)

    def do_EOF(self, _):
        """logout and quit emulator"""
        self.clt.stop()
        return True

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Game client emulator.')
    parser.add_argument('ip', metavar='ip', type=str, default='127.0.0.1', nargs='?', help='server ip addr')
    parser.add_argument('port', metavar='port', type=int, default=9090, nargs='?', help='server port')
    args = parser.parse_args(namespace=Args())

    clt = client.Client()
    code = clt.open(args.ip, args.port)
    if code != 0:
        print 'connect to server %s:%d failed: %s' % (args.ip, args.port, os.strerror(code))
        sys.exit(0)
    print 'connect to server %s:%d success' % (args.ip, args.port)
    clt.start()
    c = CLI(clt)
    try:
        c.cmdloop()
    except Exception as ex:
        print ex
        clt.stop()
    clt.join()
