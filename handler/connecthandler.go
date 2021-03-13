package handler

import "sync"

var connonce sync.Once
var conninstance *ConnectHandler

//连接管理单例
func ConnectHandlerIns() *ConnectHandler {
	connonce.Do(func() {
		conninstance = new(ConnectHandler)
	})
	return conninstance
}

/*
连接管理
*/
type ConnectHandler struct {
	sync.Map
}

//增
func (c *ConnectHandler) C(k, v interface{}) {
	c.Store(k, v)
}

//删
func (c *ConnectHandler) D(k interface{}) {
	c.Delete(k)
}

//改
func (c *ConnectHandler) U(k, v interface{}) (interface{}, bool) {
	return c.LoadOrStore(k, v)
}

//查
func (c *ConnectHandler) R(k interface{}) (interface{}, bool) {
	return c.Load(k)
}

//遍历
func (c *ConnectHandler) Each(fu func(k, v interface{}) bool) {
	c.Range(func(k, v interface{}) bool {
		return fu(k.(interface{}), v.(interface{}))
	})
}

//大小
func (c *ConnectHandler) Size() int64 {
	var size int64
	c.Range(func(k, v interface{}) bool {
		size++
		return true
	})
	return size
}

//清空
func (c *ConnectHandler) Clear() {
	c.Range(func(key, value interface{}) bool {
		c.Delete(key)
		return true
	})
}
