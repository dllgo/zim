package connect

import (
	"errors"
	"fmt"
	"log"

	"github.com/dllgo/zim/protocol"
	"github.com/panjf2000/gnet"
)

//
type Connection struct {
	cid         string
	conn        gnet.Conn
	messagePack protocol.IMessagePack
}

//
func NewConnection(cid string, conn gnet.Conn) IConnection {
	return &Connection{
		cid:         cid,
		conn:        conn,
		messagePack: protocol.NewMessagePack(),
	}
}

//
func (c *Connection) Send(reqData []byte) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[Connection.Send]recover send data error:%v", err)
		}
	}()
	if c.conn == nil {
		return errors.New("[Connection.Send] conn not exist")
	}
	return c.conn.AsyncWrite(reqData)
}

//
func (c *Connection) Read() (protocol.Imessage, error) {
	if c.conn == nil {
		return nil, errors.New("[Read] conn not exist")
	}

	message, err := c.messagePack.UnPack(c.conn.Read())
	if err != nil {
		return nil, fmt.Errorf("Read %v", err)
	}

	return message, nil
}
