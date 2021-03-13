package handler

import (
	"fmt"
	"log"
	"sync"

	"github.com/dllgo/zim/context"
	"github.com/dllgo/zim/protocol"
	"github.com/panjf2000/gnet"
)

var workonce sync.Once
var workinstance *WorkHandler

//工作线程单例
func WorkHandlerIns() *WorkHandler {
	workonce.Do(func() {
		workinstance = &WorkHandler{
			messagePack: protocol.NewMessagePack(),
		}
	})
	return workinstance
}

/*
工作线程
*/
type WorkHandler struct {
	messagePack protocol.IMessagePack
}

/**
处理接收到的消息
*/
func (wh *WorkHandler) handleFrame(frame []byte, c gnet.Conn) {
	ctx := c.Context().(context.Context)
	connid := ctx.Value("connid").(string)
	log.Println("[TcpHandler] handle 接收到", connid, "的消息")
	// 解析收到的二进制消息
	message, err := wh.messagePack.UnPack(frame)
	if err != nil {
		log.Println(err)
	}
	// 调用用户方法
	context, err := wh.HandlerMessage(message)
	if err != nil {
		log.Println(err)
	}
	// 获取用户方法返回的结果
	result, err := context.GetResult()
	if err != nil {
		log.Println(err)
	}

	resultMessage := protocol.NewMessage([]byte("1"), result, message.GetCodecType())
	rb, err := wh.messagePack.Pack(resultMessage)
	if err != nil {
		log.Println("[TcpHandler handle] error: %v", err)
	}
	// 给客户端返回处理结果
	c.AsyncWrite(rb)
}

//
//
// /**
// 处理接收到的消息
// 处理粘包
// */
// func (wh *workHandler) handle(mp protocol.MessagePack, frame []byte, c gnet.Conn) {
//
// 	ctx := c.Context().(context.Context)
// 	connid := ctx.Value("connid").(string)
//
// 	messageCount := uint32(0)
// 	data := []byte{}
// 	for {
// 		if messageCount == 0 {
// 			data = frame
// 		} else if len(data) > int(messageCount) {
// 			log.Println(connid, "========11111==========")
// 			log.Println(connid, len(data), int(messageCount))
// 			log.Println(string(data))
// 			log.Println(connid, "==========111111========")
// 			data = data[messageCount:]
// 			log.Println(connid, "==========2222========")
// 			log.Println(string(data))
// 			log.Println(connid, "==========22222========")
// 		} else {
// 			break
// 		}
//
// 		message, err := mp.UnPack(data)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		log.Println(connid, "==============data==========")
// 		log.Println(string(message.GetData()))
// 		log.Println(connid, "==============data==========")
// 		messageCount = message.GetCount()
//
// 		context, err := gh.gmsServer.HandlerMessage(message)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		result, err := context.GetResult()
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		c.AsyncWrite(result)
// 	}
//
// }

/*
处理方法
*/
func (wh *WorkHandler) HandlerMessage(message protocol.Imessage) (*context.Context, error) {
	// log.Println(string(message.GetExt()))
	handlerfunc, ok := RouteHandlerIns().R(string(message.GetExt()))
	if !ok {
		log.Println("[HandlerMessage] Router:", message.GetExt(), " not found")
		return nil, fmt.Errorf("No Router")
	}

	// todo 可以考虑使用 pool
	context := context.NewContext()
	context.SetMessage(message)
	// 调用方法
	err := handlerfunc(context)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(" fail", err)
	}

	resultData, err := context.GetResult()
	if err != nil {
		log.Println(err)
		// todo 回写错误信息
		return nil, fmt.Errorf("", err)
	}
	log.Println(string(resultData))
	// todo 回写执行结果
	return context, nil

}
