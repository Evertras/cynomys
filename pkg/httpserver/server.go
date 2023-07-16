package httpserver

import (
	"net/http"

	"github.com/evertras/gonsen"
)

type Server struct {
	server *http.Server
}

type Config struct {
	Addr string
}

func NewServer(config Config) *Server {
	// Config sanity checks
	if config.Addr == "" {
		panic("HTTP server address not given in config")
	}

	s := http.NewServeMux()

	source := gonsen.NewSource(siteFilesRaw)

	pageIndex := gonsen.NewPage(source, "index.html", func(r *http.Request) (IndexData, int) {
		return IndexData{}, http.StatusOK
	})

	s.Handle("/", pageIndex)

	return &Server{
		server: &http.Server{
			Addr:    config.Addr,
			Handler: s,
		},
	}
}

func (s *Server) ServeAndListen() error {
	return s.server.ListenAndServe()
}
