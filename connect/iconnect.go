package connect

import "github.com/dllgo/zim/protocol"

type IConnection interface {
	Send(reqData []byte) error
	Read() (protocol.Imessage, error)
}
