package handler

import (
	"fmt" 
	"log" 
	"context"
	"time"
	"strconv"

	"github.com/panjf2000/gnet" 
	"github.com/panjf2000/gnet/pool/goroutine" 
)

type WSSHandler struct {
	*gnet.EventServer
	codec       gnet.ICodec
	pool        *goroutine.Pool
	gnetServer  gnet.Server 
}

func NewWSSHandler(srv  IServer) Ihandler {
	return &WSSHandler{
		pool:  			goroutine.Default(), 
	} 
}
/*
回收资源
*/
func (ws *WSSHandler) Release() {
	log.Println("[WSSHandler] stop")
	ws.pool.Release() 
}
/**
处理接收到的消息
*/
func (ws *WSSHandler) handle(frame []byte, c gnet.Conn) {
	log.Println("[WSSHandler] handle")
	// 给客户端返回处理结果
	c.AsyncWrite(frame)
}