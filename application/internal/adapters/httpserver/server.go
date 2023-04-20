package httpserver

// @title           Bakery backend API
// @version         1.0
// @description     This is a bakery backend server

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
import (
	"context"
	"errors"
	"net/http"
	"time"

	"bakery/application/docs"
	"bakery/application/internal/config"
	"bakery/application/internal/domain/bakery"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	Server *http.Server
	Bakery *bakery.Bakery
	Log    *logrus.Logger
}
type route struct {
	Method  string
	Path    string
	Handler func(s *Server) http.HandlerFunc
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
	docs.SwaggerInfo.BasePath = "/"
	return &server
}

func (s *Server) addRoutes() http.Handler {
	s.Log.Info("Routes", routes)
	router := *chi.NewRouter()
	router.Use(s.loggingMiddleware)
	for _, mRoute := range routes {
		router.Options(mRoute.Path, mRoute.Handler(s))
		router.Method(mRoute.Method, mRoute.Path, mRoute.Handler(s))
	}
	router.Method("GET", "/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	return &router
}

func (s *Server) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		shutDownCtx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
		defer cancelFunc()
		s.Server.Shutdown(shutDownCtx)
	}()
	err := s.Server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		s.Log.Error(err)
		s.Log.Fatal("Error starting http server")
	} else {
		s.Log.Log(logrus.DebugLevel, "Http server stopped")
	}

}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		s.Log.Info(r.RequestURI, " ", r.Method)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
