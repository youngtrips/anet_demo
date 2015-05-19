package anet

type Message struct {
	Api     int16
	Payload interface{}
}

func NewMessage(api int16, payload interface{}) *Message {
	return &Message{
		Api:     api,
		Payload: payload,
	}
}
