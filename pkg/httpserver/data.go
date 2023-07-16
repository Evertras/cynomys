package httpserver

type ListenerInfo struct {
	Address string
}

type IndexData struct {
	ListenersTCP []ListenerInfo
	ListenersUDP []ListenerInfo
}
