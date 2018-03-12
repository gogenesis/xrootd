package read

const RequestID uint16 = 3013

type Response struct {
	Data []byte
}

type Request struct {
	FileHandle [4]byte
	Offset     int64
	Length     int32
	ArgsLength int32
}

func NewRequest(fileHandle [4]byte, offset int64, length int32) Request {
	return Request{fileHandle, offset, length, 0}
}
