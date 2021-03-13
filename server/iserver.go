package server

import (
	"github.com/dllgo/zim/context"
)

//
type IServer interface {
	// 启动zim服务
	Serve(ip string, port int)
	// 停止zim服务
	Stop()
	// 注册处理器
	AddRouter(handlerName string, handlerFunc context.HandlerFunc)
	// 获取处理器
	GetRouter(handlerName string) (context.HandlerFunc, bool)
}
