package dirlist

const RequestID uint16 = 3004

type Response struct {
	Data []byte
}

type Request struct {
	Reserved1  [15]byte
	Options    byte
	PathLength int32
	Path       []byte
}

func NewRequest(path string) Request {
	var pathBytes = make([]byte, len(path))
	copy(pathBytes, path)

	return Request{[15]byte{}, 0, int32(len(path)), pathBytes}
}
