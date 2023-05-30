package httpserver

import (
	"io/fs"
	"net/http"
)

type Server struct {
	server *http.Server
}

type Config struct {
	Addr string
}

func NewServer(config Config) *Server {
	s := http.NewServeMux()

	siteFiles, err := fs.Sub(siteFilesRaw, "site")

	if err != nil {
		panic(err)
	}

	s.Handle("/", http.FileServer(http.FS(siteFiles)))

	if config.Addr == "" {
		panic("HTTP server address not given in config")
	}

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
