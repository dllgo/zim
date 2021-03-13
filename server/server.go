package server

import (
	"sync"

	"github.com/dllgo/zim/context"
	"github.com/dllgo/zim/protocol"
	"github.com/dllgo/zim/handler"
)


type server struct {
	// 整个服务级别的锁
	sync.RWMutex
	// 路由Map
	routerMap map[string]context.HandlerFunc
	// tcp 服务
	tcpHandler handler.Ihandler 
	// wss 服务
	wssHandler handler.Ihandler 
}

/*
初始化GMS服务
*/
func NewServer() IServer {
	s := server{
		routerMap: make(map[string]context.HandlerFunc),
	}
	return &s
}

/*
准备启动服务的资源
*/
func (s *server) InitServe(port int) {
	log.Println("[Server] InitServe")
	// 启动gnet
	s.tcpHandler = NewTcpHandler(s)
	s.wssHandler = NewWSSHandler(s)
	if port < 1 {
		port = 9000
	}
	log.Fatal(gnet.Serve(
		s.evHandler,
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

	// 准备启动服务的资源
	s.InitServe(port)

}

/*
停止服务 回收资源
*/
func (s *server) Stop() {
	log.Println("[server] stop")
	s.evHandler.Release()
}

/*
添加路由
*/
func (s *server) AddRouter(handlerName string, handlerFunc context.HandlerFunc) {
	s.Lock()
	defer s.Unlock()

	// 注册路由
	if _, ok := s.routerMap[handlerName]; ok {
		log.Println("[AddRouter] fail handlerName:", handlerName, " alread exist")
		return
	}
	s.routerMap[handlerName] = handlerFunc
}

/*
获取路由
*/
func (s *server) GetRouter(handlerName string) (context.HandlerFunc, error) {
	s.RLock()
	defer s.RUnlock()
	if handler, ok := s.routerMap[handlerName]; ok {
		return handler, nil
	}
	return nil, errors.New("[GetRouter] Router not found")
}

/*
处理方法
*/
// func (s *server) HandlerMessage(message protocol.Imessage) (*gmsContext.Context, error) {
// func (s *server) HandlerMessage(message protocol.Imessage) (gmsContext.Context, error) {
func (s *server) HandlerMessage(message protocol.Imessage) (*context.Context, error) {
	// log.Println(string(message.GetExt()))
	handler, err := s.GetRouter(string(message.GetExt()))
	if err != nil {
		log.Println("[HandlerMessage] Router:", message.GetExt(), " not found", err)
		return nil, fmt.Errorf("No Router", err)
	}

	// todo 可以考虑使用 pool
	context := context.NewContext()
	context.SetMessage(message)
	// 调用方法
	err = handler(context)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(" fail", err)
	}

	// resultData, err := context.GetResult()
	// if err != nil {
	// 	log.Println(err)
	// 	// todo 回写错误信息
	// 	return nil, fmt.Errorf("", err)
	// }
	// log.Println(string(resultData))
	// todo 回写执行结果
	return context, nil

}