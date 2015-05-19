package anet

type Protocol interface {
	Encode(api int16, data interface{}) ([]byte, error)
	Decode(data []byte) (int16, interface{}, error)
}
