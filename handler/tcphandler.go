package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/dllgo/zim/utils"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
)

var tcponce sync.Once
var tcpinstance *TCPHandler

//Tcp 服务单例
func TCPHandlerIns() *TCPHandler {
	tcponce.Do(func() {
		tcpinstance = &TCPHandler{
			pool: goroutine.Default(),
		}
	})
	return tcpinstance
}

//tcp event
type TCPHandler struct {
	*gnet.EventServer
	codec      gnet.ICodec
	pool       *goroutine.Pool
	gnetServer gnet.Server
}

/*
回收资源
*/
func (eh *TCPHandler) Release() {
	log.Println("[TcpHandler] stop")
	eh.pool.Release()
}

/*
gnet 服务启动成功
*/
func (eh *TCPHandler) OnInitComplete(server gnet.Server) (action gnet.Action) {
	log.Printf("[TcpHandler OnInitComplete] listening on %s (multi-cores: %t, loops: %d)\n",
		server.Addr.String(), server.Multicore, server.NumEventLoop)
	eh.gnetServer = server
	return
}

/*
gnet 新建连接
*/
func (eh *TCPHandler) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	connid := utils.GenCid()
	ctx := context.WithValue(context.Background(), "cid", connid)
	log.Println(fmt.Sprintf("[TcpHandler OnOpened] client: %v open. RemoteAddr:%v", connid, c.RemoteAddr().String()))
	log.Println("[TcpHandler OnOpened] Conn count:", eh.gnetServer.CountConnections())
	ConnectHandlerIns().C(connid, c)
	c.SetContext(ctx)
	return
}

/*
gnet 连接断开
*/
func (eh *TCPHandler) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		log.Println("[TcpHandler OnClosed] error:", err)
		return
	}
	ctx := c.Context().(context.Context)
	cid := ctx.Value("cid").(string)
	ConnectHandlerIns().D(cid)
	log.Println("[TcpHandler OnClosed] client: " + utils.GetAddrByCid(cid) + " Close;===Conn count:" + strconv.FormatInt(eh.Size(), 10))
	return
}

// 定时器
func (eh *TCPHandler) Tick() (delay time.Duration, action gnet.Action) {
	log.Println("[TcpHandler OnClosed] Tick: " + strconv.FormatInt(eh.Size(), 10))
	ConnectHandlerIns().Each(func(key, value interface{}) bool {
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

//接收数据
func (eh *TCPHandler) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// Use ants pool to unblock the event-loop.
	err := eh.pool.Submit(func() {
		WorkHandlerIns().handleFrame(frame, c)
	})

	if err != nil {
		log.Println("[React] error:", err)
	}
	return
}

// Size 在线人数
func (eh *TCPHandler) Size() int64 {
	return ConnectHandlerIns().Size()

}
