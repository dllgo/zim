package zim
import (
	"sync"

	"github.com/dllgo/zim/context"
	"github.com/dllgo/zim/server"
)


type zim struct {
	server server.IServer
}

var defaultZim = newZim()

/*
初始化zim
*/
func newZim() *zim {
	zim := zim{
		server: server.NewServer(),
	}
	return &zim
}
/*
启动zim
*/
func Serve(ip string, port int) {
	// 启动GMS服务
	defaultZim.server.Serve(ip, port)
}
/**
添加服务路由
*/
func AddRouter(handlerName string, handlerFunc context.HandlerFunc) {
	defaultZim.server.AddRouter(handlerName, handlerFunc)
}
