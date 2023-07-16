package httpserver

import (
	"net/http"

	"github.com/evertras/cynomys/pkg/listener"

	"github.com/evertras/gonsen"
)

type Server struct {
	server *http.Server
}

type OverallStatusGetter interface {
	TCPListeners() []*listener.TCPListener
	UDPListeners() []*listener.UDPListener
}

func NewServer(addr string, statusGetter OverallStatusGetter) *Server {
	// Config sanity checks
	if addr == "" {
		// This might be from bad user config, so use friendlier message
		panic("http server address not given")
	}
	if statusGetter == nil {
		// This is from a bad code path
		panic("statusGetter isn't set")
	}

	s := http.NewServeMux()

	source := gonsen.NewSource(siteFilesRaw)

	pageIndex := gonsen.NewPage(source, "index.html", func(r *http.Request) (IndexData, int) {
		return IndexData{overallStatusFromGetter(statusGetter)}, http.StatusOK
	})

	s.Handle("/", pageIndex)

	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: s,
		},
	}
}

func (s *Server) ServeAndListen() error {
	return s.server.ListenAndServe()
}
