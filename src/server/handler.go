package server

import (
	"anet"
	"db"
	"log"
	"protocol"
)

func (app *App) onEvent(ev anet.Event) {
	switch ev.Type {
	case anet.EVENT_ACCEPT:
		log.Printf("new connection...")
		session := ev.Session
		if session != nil {
			session.Start(app.events)
		}
		break
	case anet.EVENT_CONNECT_SUCCESS:
		break
	case anet.EVENT_CONNECT_FAILED:
		break
	case anet.EVENT_DISCONNECT:
		log.Printf("connection closed...")
		onLogut(app, ev.Session)
		break
	case anet.EVENT_MESSAGE:
		msg := ev.Data.(*anet.Message)
		app.onMessage(ev.Session, msg.Api, msg.Payload)
		break
	case anet.EVENT_RECV_ERROR:
		break
	case anet.EVENT_SEND_ERROR:
		break
	default:
		log.Printf("invalid event type: %d", ev.Type)
		break
	}
}

func (app *App) onMessage(session *anet.Session, opcode int16, payload interface{}) {
	log.Printf("opcode: %d, payload: %v", opcode, payload)
	switch protocol.OPCODE(opcode) {
	case protocol.OPCODE_LOGIN_REQ:
		msg := payload.(*protocol.LoginReq)
		onLogin(app, session, msg)
	case protocol.OPCODE_ONLINE_USERS_REQ:
		msg := payload.(*protocol.NullMessage)
		onReqOnlineUsers(app, session, msg)
		break
	case protocol.OPCODE_EXEC_CMD_REQ:
		msg := payload.(*protocol.ExecCmdReq)
		onExecCmdReq(app, session, msg)
	}
}

func Send(session *anet.Session, opcode protocol.OPCODE, msg interface{}) {
	session.Send(int16(opcode), msg)
}

func CommonAck(session *anet.Session, opcode protocol.OPCODE, errno protocol.ERROR) {
	ack := protocol.CommonAck{
		Opcode: opcode.Enum(),
		Errno:  errno.Enum(),
	}
	Send(session, protocol.OPCODE_COMMON_ACK, &ack)
}

func onLogin(app *App, session *anet.Session, msg *protocol.LoginReq) {
	log.Printf("on login: username[%s], password[%s]", msg.GetUsername(), msg.GetPassword())
	ret, uid := db.AccountAuth(msg.GetUsername(), msg.GetPassword())
	if ret != protocol.ERROR_SUCCESS {
		CommonAck(session, protocol.OPCODE_LOGIN_REQ, ret)
	} else {
		if _, present := app.users[uid]; present {
			CommonAck(session, protocol.OPCODE_LOGIN_REQ, protocol.ERROR_ALREADY_LOGINED)
			return
		}
		log.Printf("ret=%d, uid=%d", ret, uid)
		user := db.LoadUser(uid)
		log.Printf("%v", user)
		info := protocol.UserInfo{
			Id:       &user.Id,
			Username: &user.Name,
		}
		ack := protocol.LoginAck{
			Info: &info,
		}
		Send(session, protocol.OPCODE_LOGIN_ACK, &ack)

		usession := &UserSession{
			Id:      user.Id,
			User:    user,
			Session: session,
		}
		//
		app.users[usession.Id] = usession
		app.sessionMapping[session.ID()] = user.Id
	}
}

func onLogut(app *App, session *anet.Session) {
	if uid, present := app.sessionMapping[session.ID()]; present {
		delete(app.users, uid)
		delete(app.sessionMapping, session.ID())
	}
}

func onReqOnlineUsers(app *App, session *anet.Session, msg *protocol.NullMessage) {
	infos := make([]*protocol.UserInfo, 0)
	for _, usession := range app.users {
		info := protocol.UserInfo{
			Id:       &usession.User.Id,
			Username: &usession.User.Name,
		}
		infos = append(infos, &info)
	}
	ack := protocol.OnlineUserList{
		Users: infos,
	}
	Send(session, protocol.OPCODE_ONLINE_USERS_ACK, &ack)
}

func onExecCmdReq(app *App, session *anet.Session, msg *protocol.ExecCmdReq) {
	usession, present := app.users[msg.GetTargetUid()]
	ret := protocol.ERROR_SUCCESS
	if !present {
		ret = protocol.ERROR_NO_FOUND_USER
	} else {
		cmd := msg.GetCmd()
		info := protocol.ExecCmdInfo{
			Cmd: &cmd,
		}
		Send(usession.Session, protocol.OPCODE_EXEC_CMD_NTF, &info)
	}
	CommonAck(session, protocol.OPCODE_EXEC_CMD_REQ, ret)
}
