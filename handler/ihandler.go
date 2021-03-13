package handler

type Ihandler interface {
	Handle(frame []byte, c gnet.Conn)
	Release()
}