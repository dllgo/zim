package context

import "github.com/dllgo/zim/protocol"

type HandlerFunc func(c *Context) error

type IContext interface {
	// 设置参数数据
	SetMessage(message protocol.Imessage) error
	// 参数转换
	Param(interface{}) error
	// 返回参数
	Result(interface{}) error
	// 获取结果
	GetResult() ([]byte, error)
}