package httpserver

type ListenerInfo struct {
	Address string
}

type OverallStatus struct {
	TCPListeners []ListenerInfo
	UDPListeners []ListenerInfo
}

type IndexData struct {
	OverallStatus
}

func overallStatusFromGetter(getter OverallStatusGetter) OverallStatus {
	s := OverallStatus{}

	for _, tcpListener := range getter.TCPListeners() {
		s.TCPListeners = append(s.TCPListeners, ListenerInfo{
			Address: tcpListener.Addr(),
		})
	}

	for _, udpListener := range getter.UDPListeners() {
		s.UDPListeners = append(s.UDPListeners, ListenerInfo{
			Address: udpListener.Addr(),
		})
	}

	return s
}
