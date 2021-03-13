package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/dllgo/zim/codec"
)
 
/*
MessagePack 消息编码、解码
实现gnet.ICodec 接口
*/
type MessagePack struct {
	// Message
}

func NewMessagePack() IMessagePack {
	return &MessagePack{}
}

/*
Pack 消息编码
消息格式
扩展数据长度|主体数据长度|扩展数据|主体数据
*/
func (m *MessagePack) Pack(message Imessage) ([]byte, error) {
	result := make([]byte, 0)

	buffer := bytes.NewBuffer(result)

	if err := binary.Write(buffer, binary.BigEndian, message.GetExtLen()); err != nil {
		s := fmt.Sprintf("[Pack] Pack extLen error , %v", err)
		return nil, errors.New(s)
	}

	if err := binary.Write(buffer, binary.BigEndian, message.GetDataLen()); err != nil {
		s := fmt.Sprintf("[Pack] Pack dataLen error , %v", err)
		return nil, errors.New(s)
	}

	if err := binary.Write(buffer, binary.BigEndian, message.GetCodecType()); err != nil {
		s := fmt.Sprintf("[Pack] Pack GetCodecType error , %v", err)
		return nil, errors.New(s)
	}

	if message.GetExtLen() > 0 {
		if err := binary.Write(buffer, binary.BigEndian, message.GetExt()); err != nil {
			s := fmt.Sprintf("[Pack] Pack ext error , %v", err)
			return nil, errors.New(s)
		}
	}

	if message.GetDataLen() > 0 {
		if err := binary.Write(buffer, binary.BigEndian, message.GetData()); err != nil {
			s := fmt.Sprintf("[Pack] Pack data error , %v", err)
			return nil, errors.New(s)
		}
	}

	// log.Println(string(buffer.Bytes()))
	return buffer.Bytes(), nil
}

/*
UnPack 消息解码
消息格式
扩展数据长度|主体数据长度|扩展数据|主体数据
*/
func (m *MessagePack) UnPack(binaryMessage []byte) (Imessage, error) {
	// log.Println("1:binaryMessage:", string(binaryMessage))
	header := bytes.NewReader(binaryMessage[:HeaderLength])

	// 只解压head的信息，得到dataLen和msgID
	var extLen, dataLen uint32
	// 读取扩展信息长度
	if err := binary.Read(header, binary.BigEndian, &extLen); err != nil {
		return nil, err
	}

	// 读入消息长度
	if err := binary.Read(header, binary.BigEndian, &dataLen); err != nil {
		return nil, err
	}

	msg := &Message{
		extLen:  extLen,
		dataLen: dataLen,
		count:   HeaderLength + extLen + dataLen + 1,
	}

	codecType := binaryMessage[HeaderLength : HeaderLength+1]
	msg.codecType = codec.CodecType(codecType[0])

	// 截取消息头后的所有内容
	content := binaryMessage[HeaderLength+1 : msg.GetCount()]

	// 获取扩展消息
	msg.SetExt(content[:msg.GetExtLen()])
	// 获取消息内容
	msg.SetData(content[msg.GetExtLen():])
	return msg, nil
}

func (m *MessagePack) ReadUnPack(conn net.Conn) (Imessage, error) {
	headData := make([]byte, HeaderLength)
	_, err := io.ReadFull(conn, headData) // ReadFull 会把msg填充满为止
	if err != nil {
		log.Println("[Read] read header error", err)
		return nil, err
	}

	header := bytes.NewReader(headData)
	// log.Println(string(headData))

	// 只解压head的信息，得到dataLen和msgID
	var extLen, dataLen uint32
	// 读取扩展信息长度
	if err := binary.Read(header, binary.BigEndian, &extLen); err != nil {
		return nil, err
	}

	// 读入消息长度
	if err := binary.Read(header, binary.BigEndian, &dataLen); err != nil {
		return nil, err
	}

	msg := &Message{
		extLen:  extLen,
		dataLen: dataLen,
		count:   HeaderLength + extLen + dataLen + 1,
	}

	codecType := make([]byte, 1)
	{
		n, err := io.ReadFull(conn, codecType)
		if err != nil {
			log.Println("[Read] read codecType error", err)
			return nil, err
		}
		if uint32(n) != 1 {
			log.Println("[Read] read codecType len error")
			return nil, errors.New("read codecType error")
		}
	}

	extData := make([]byte, extLen)
	// 读取扩展信息
	{
		n, err := io.ReadFull(conn, extData)
		if err != nil {
			log.Println("[Read] read extData error", err)
			return nil, err
		}
		if uint32(n) != extLen {
			log.Println("[Read] read extData len error")
			return nil, errors.New("read extData error")
		}
	}

	data := make([]byte, dataLen)
	// 读取数据
	{
		n, err := io.ReadFull(conn, data)
		if err != nil {
			log.Println("[Read] read date error", err)
			return nil, err
		}
		if uint32(n) != dataLen {
			log.Println("[Read] read data len error")
			return nil, errors.New("read extData error")
		}
	}

	// 设置编码类型
	msg.SetCodecType(codec.CodecType(codecType[0]))
	// 获取扩展消息
	msg.SetExt(extData)
	// 获取消息内容
	msg.SetData(data)
	return msg, nil

}