package handler

import (
	"log"
	"sync"

	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
)

var wssonce sync.Once
var wssinstance *WSSHandler

//wss 服务单例
func WSSHandlerIns() *WSSHandler {
	wssonce.Do(func() {
		wssinstance = &WSSHandler{
			pool: goroutine.Default(),
		}
	})
	return wssinstance
}

//wss event
type WSSHandler struct {
	*gnet.EventServer
	codec      gnet.ICodec
	pool       *goroutine.Pool
	gnetServer gnet.Server
}

/*
回收资源
*/
func (ws *WSSHandler) Release() {
	log.Println("[wssHandler] stop")
	ws.pool.Release()
}

/**
处理接收到的消息
*/
func (ws *WSSHandler) handle(frame []byte, c gnet.Conn) {
	log.Println("[wssHandler] handle")
	// 给客户端返回处理结果
	c.AsyncWrite(frame)
}
