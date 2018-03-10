package protocol

const RequestID uint16 = 3006

type Response struct {
	BinaryProtocolVersion int32
	Flags                 int32
}

type SecurityInfo struct {
	SecurityVersion       byte
	SecurityOptions       byte
	SecurityLevel         byte
	SecurityOverridesSize byte
}

type Request struct {
	ClientProtocolVersion int32
	Reserved1             [11]byte
	Options               byte
	Reserved2             int32
}

func NewRequest(protocolVersion int32) Request {
	return Request{protocolVersion, [11]byte{}, 0, 0}
}
