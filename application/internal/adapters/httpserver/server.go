package httpserver

import (
	"context"
	"errors"
	"net/http"

	"bakery/application/internal/config"
	"bakery/application/internal/domain/bakery"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Server *http.Server
	Bakery *bakery.Bakery
	Log    *logrus.Logger
}
type route struct {
	Methods []string
	Path    string
	Handler func(s *Server) http.Handler
}

var routes []route

func New(cfg *config.Config, b *bakery.Bakery, log *logrus.Logger) *Server {
	server := Server{
		Bakery: b,
		Log:    log,
	}

	httpServer := http.Server{
		Addr:    cfg.Server_addr + ":" + cfg.Server_port,
		Handler: server.addRoutes(),
	}
	server.Server = &httpServer
	//docs.SwaggerInfo.BasePath = "/"
	return &server
}

func (s *Server) addRoutes() http.Handler {
	router := *mux.NewRouter()
	for _, mRoute := range routes {
		router.Handle(mRoute.Path, mRoute.Handler(s)).Methods(mRoute.Methods...)
	}
	return &router
}

func (s *Server) Start() {
	err := s.Server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		s.Log.Error(err)
		s.Log.Fatal("Error starting http server")
	} else {
		s.Log.Log(logrus.DebugLevel, "Http server stopped")
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
