package handler

import (
	"fmt" 
	"log" 
	"context"
	"time"
	"strconv"

	"github.com/panjf2000/gnet" 
	"github.com/panjf2000/gnet/pool/goroutine"
	"github.com/dllgo/zim/utils"
	"github.com/dllgo/zim/protocol"
	"github.com/dllgo/zim/server"
)

type TcpHandler struct {
	*gnet.EventServer
	codec       gnet.ICodec
	pool        *goroutine.Pool
	gnetServer  gnet.Server 
	messagePack protocol.IMessagePack
	zimServer   server.IServer
}

func NewTcpHandler(srv  server.IServer) Ihandler {
	return &TcpHandler{
		pool:  			goroutine.Default(), 
		messagePack:    protocol.NewMessagePack(), 
		zimServer:		srv,
	} 
}
/*
回收资源
*/
func (eh *TcpHandler) Release() {
	log.Println("[TcpHandler] stop")
	eh.pool.Release()  
}

func (eh *TcpHandler) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// Use ants pool to unblock the event-loop.
	err := eh.pool.Submit(func() {
		eh.handle(frame, c)
	})

	if err != nil {
		log.Println("[React] error:", err)
	}
	return
}

/**
处理接收到的消息
*/
func (eh *TcpHandler) handle(frame []byte, c gnet.Conn) {
	log.Println("[TcpHandler] handle")
	// 解析收到的二进制消息
	message, err := eh.messagePack.UnPack(frame)
	if err != nil {
		log.Println(err)
	}
	// 调用用户方法
	context, err := eh.gmsServer.HandlerMessage(message)
	if err != nil {
		log.Println(err)
	}
	// 获取用户方法返回的结果
	result, err := context.GetResult()
	if err != nil {
		log.Println(err)
	}

	resultMessage := NewMessage([]byte("1"), result, message.GetCodecType())
	rb, err := eh.messagePack.Pack(resultMessage)
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
// func (gh *eventHandler) handle(mp protocol.MessagePack, frame []byte, c gnet.Conn) {
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
gnet 服务启动成功
*/
func (eh *TcpHandler) OnInitComplete(server gnet.Server) (action gnet.Action) {
	log.Printf("[TcpHandler OnInitComplete] listening on %s (multi-cores: %t, loops: %d)\n",
	server.Addr.String(), server.Multicore, server.NumEventLoop)
	eh.gnetServer = server
	return
}

/*
gnet 新建连接
*/
func (eh *TcpHandler) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	connid := utils.GenCid()
	ctx := context.WithValue(context.Background(), "cid", connid)
	log.Println(fmt.Sprintf("[TcpHandler OnOpened] client: %v open. RemoteAddr:%v", connid, c.RemoteAddr().String()))
	log.Println("[TcpHandler OnOpened] Conn count:", eh.gnetServer.CountConnections())
	utils.SyncMapIns().C(connid,c)
	c.SetContext(ctx)
	return
}

/*
gnet 连接断开
*/
func (eh *TcpHandler) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		log.Println("[TcpHandler OnClosed] error:", err)
		return
	}
	ctx := c.Context().(context.Context)
	cid := ctx.Value("cid").(string)
	utils.SyncMapIns().D(cid)
	log.Println("[TcpHandler OnClosed] client: " + utils.GetAddrByCid(cid) + " Close;===Conn count:"+ strconv.FormatInt(eh.Size(),10))
	return
}

// 定时器
func (eh *TcpHandler) Tick() (delay time.Duration, action gnet.Action) { 
	log.Println("[TcpHandler OnClosed] Tick: "+ strconv.FormatInt(eh.Size(),10))
	utils.SyncMapIns().Each(func(key, value interface{}) bool {
		addr := key.(string)
		c := value.(gnet.Conn)
		c.AsyncWrite([]byte(fmt.Sprintf("heart beating to %s\n", addr)))
		return true
	})
	var interval time.Duration
	interval = time.Second
	delay = interval
	return
}

// Size 在线人数
func (eh *TcpHandler)Size() int64 { 
	return utils.SyncMapIns().Size()

}