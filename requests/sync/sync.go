package sync

const RequestID uint16 = 3016

type Request struct {
	FileHandle [4]byte
	Reserved1  [12]byte
	Reserved2  int32
}

func NewRequest(fileHandle [4]byte) Request {
	return Request{fileHandle, [12]byte{}, 0}
}
