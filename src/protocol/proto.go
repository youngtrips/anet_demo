package protocol

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"reflect"
)

var (
	protoMapping map[OPCODE]reflect.Type
)

func init() {
	protoMapping = make(map[OPCODE]reflect.Type)

	Register(OPCODE_LOGIN_REQ, LoginReq{})
	Register(OPCODE_ONLINE_USERS_REQ, NullMessage{})
	Register(OPCODE_EXEC_CMD_REQ, ExecCmdReq{})
}

type Proto struct {
}

func (self Proto) Encode(opcode int16, msg interface{}) ([]byte, error) {
	buf := make([]byte, 2)
	buf[0] = byte(opcode >> 8)
	buf[1] = byte(opcode & 0xFF)

	data, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		return nil, err
	}
	buf = append(buf, data...)
	return buf, nil
}

func (self Proto) Decode(data []byte) (int16, interface{}, error) {
	opcode := int16(0)
	opcode |= int16(data[0]) << 8
	opcode |= int16(data[1]) & 0xFF

	typeValue, present := protoMapping[OPCODE(opcode)]
	if !present {
		return 0, nil, errors.New("no such protocal type")
	}

	value := reflect.New(typeValue)
	msg := value.Interface().(proto.Message)
	if err := proto.Unmarshal(data[2:], msg); err != nil {
		return 0, nil, err
	}
	return opcode, msg, nil
}

func Register(op OPCODE, proto interface{}) {
	protoMapping[op] = reflect.TypeOf(proto)
}
