package handler

import (
	"sync"

	"github.com/dllgo/zim/context"
)

var routeonce sync.Once
var routeinstance *RouterHandler

//路由单例
func RouteHandlerIns() *RouterHandler {
	routeonce.Do(func() {
		routeinstance = new(RouterHandler)
	})
	return routeinstance
}

/*
路由
*/
type RouterHandler struct {
	sync.Map
}

//增
func (s *RouterHandler) C(handlerName string, handlerFunc context.HandlerFunc) {
	s.Store(handlerName, handlerFunc)
}

//删
func (s *RouterHandler) D(handlerName string) {
	s.Delete(handlerName)
}

//查
func (s *RouterHandler) R(handlerName string) (context.HandlerFunc, bool) {
	if v, ok := s.Load(handlerName); ok {
		return v.(context.HandlerFunc), ok
	}
	return nil, false
}

//遍历
func (s *RouterHandler) Each(fu func(handlerName string, handlerFunc context.HandlerFunc) bool) {
	s.Range(func(k, v interface{}) bool {
		return fu(k.(string), v.(context.HandlerFunc))
	})
}

//大小
func (s *RouterHandler) Size() int64 {
	var size int64
	s.Range(func(k, v interface{}) bool {
		size++
		return true
	})
	return size
}

//清空
func (s *RouterHandler) Clear() {
	s.Range(func(key, value interface{}) bool {
		s.Delete(key)
		return true
	})
}
