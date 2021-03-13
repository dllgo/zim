package utils

import (
	"sync"
)
var once sync.Once
var instance *SyncMap
func SyncMapIns() *SyncMap {
    once.Do(func() {
        instance = new(SyncMap)
    })
    return instance
}

type SyncMap struct {
	sync.Map
}
  
//增
func (s *SyncMap) C(k, v interface{}) {
	s.Store(k,v)
}
//删
func (s *SyncMap) D(k interface{}){
	s.Delete(k)
}
//改
func (s *SyncMap) U(k, v interface{}) (interface{}, bool){
	return s.LoadOrStore(k,v)
}
//查
func (s *SyncMap) R(k interface{}) (interface{}, bool){
	return s.Load(k)
}
//遍历
func (s *SyncMap) Each(fu func(k, v interface{}) bool) {
	s.Range(func(k, v interface{}) bool {
		return fu(k.(interface{}), v.(interface{}))
	})
} 
//大小
func (s *SyncMap) Size() int64 {
	var size int64
	s.Range(func(k, v interface{}) bool {
		size++
		return true
	})
	return size
}
//清空
func (s *SyncMap) Clear() {
	s.Range(func(key, value interface{}) bool {
		s.Delete(key)
		return true
	})
}
