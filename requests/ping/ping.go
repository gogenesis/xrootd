package ping

const RequestID uint16 = 3011

type Request struct {
	Reserved1 [16]byte
	Reserved2 int32
}

func NewRequest() Request {
	return Request{[16]byte{}, 0}
}
