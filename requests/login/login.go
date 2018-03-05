package login

import (
	"os"
)

const RequestID uint16 = 3007

type Response struct {
	SecurityInformation []byte
}

type Request struct {
	Pid           int32
	UsernameBytes [8]byte
	Reserved1     byte
	Ability       byte
	Reserved2     byte
	Reserved3     byte
	Reserved4     int32
}

func NewRequest(username string) Request {
	var usernameBytes [8]byte
	copy(usernameBytes[:], username)
	var ability = byte(0)

	return Request{int32(os.Getpid()), usernameBytes, 0, ability, 0, 0, 0}
}
