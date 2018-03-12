package bind

const RequestID uint16 = 3024

type Response struct {
	PathID byte
}

type Request struct {
	SessionID [16]byte
	Reserved  int32
}

func NewRequest(sessionID [16]byte) Request {
	return Request{sessionID, 0}
}
