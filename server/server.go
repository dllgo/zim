package server

import (
	"fmt"
	"log"
	"time"

	"github.com/dllgo/zim/context"
	"github.com/dllgo/zim/handler"
	"github.com/panjf2000/gnet"
)

//
type server struct {
}

/*
初始化zim服务
*/
func NewServer() IServer {
	s := server{}
	return &s
}

/*
准备启动服务的资源
*/
func (s *server) StartTcpServe(port int) {
	log.Println("[ZIMServer] StartTcpServe")
	// 启动tcp
	if port < 1 {
		port = 9000
	}
	log.Fatal(gnet.Serve(
		handler.TCPHandlerIns(),
		fmt.Sprintf("tcp://:%v", port),
		gnet.WithMulticore(true),
		gnet.WithTCPKeepAlive(time.Minute*5), // todo 需要确定是否对长连接有影响
		gnet.WithTicker(true),
	))
}

/*
准备启动服务的资源
*/
func (s *server) StartWSSServe(port int) {
	log.Println("[ZIMServer] StartWSSServe")
	// 启动websocket
	if port < 1 {
		port = 9001
	}
	log.Fatal(gnet.Serve(
		handler.WSSHandlerIns(),
		fmt.Sprintf("tcp://:%v", port),
		gnet.WithMulticore(true),
		gnet.WithTCPKeepAlive(time.Minute*5), // todo 需要确定是否对长连接有影响
		gnet.WithTicker(true),
	))
}

/*
启动服务
*/
func (s *server) Serve(ip string, port int) {
	log.Println("[server] start run Server")
	// // 启动所有插件
	// s.plugins.Start()

	// // 注册插件
	// s.plugins.Registe(ip, port)
	go s.StartWSSServe(port + 1)
	// 准备启动服务的资源
	s.StartTcpServe(port)

}

/*
停止服务 回收资源
*/
func (s *server) Stop() {
	log.Println("[server] stop")

}

/*
添加路由
*/
func (s *server) AddRouter(handlerName string, handlerFunc context.HandlerFunc) {
	handler.RouteHandlerIns().C(handlerName, handlerFunc)
}

/*
获取路由
*/
func (s *server) GetRouter(handlerName string) (context.HandlerFunc, bool) {
	return handler.RouteHandlerIns().R(handlerName)

}
