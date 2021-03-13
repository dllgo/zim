package context

import (
	"context"
	"log"

	"github.com/dllgo/zim/codec"
	"github.com/dllgo/zim/protocol"
)


type Context struct {
	context.Context // todo context 功能待完善（参考gin的context实现）
	message         protocol.Imessage
	resultData      []byte
}
 
func NewContext() *Context {
	return &Context{}
}

func (c *Context) SetMessage(message protocol.Imessage) error {
	c.message = message
	return nil
}

/**
把请求中的信息反序列化成用户指定的对象
*/
func (c *Context) Param(param interface{}) error {

	// 获取指定的序列化器
	codec := codec.GetCodec(c.message.GetCodecType())
	err := codec.Decode(c.message.GetData(), param)

	// err := json.Unmarshal(c.message.GetData(), param)
	if err != nil {
		log.Println("[Param] error", err)
		return err
	}
	return nil
}

func (c *Context) Result(result interface{}) error {
	codec := codec.GetCodec(c.message.GetCodecType())
	r, err := codec.Encode(result)
	// // todo 改为其他序列化方式
	// r, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	// log.Println(len(r))
	c.resultData = r
	return nil
}

func (c *Context) GetResult() ([]byte, error) {
	return c.resultData, nil
}