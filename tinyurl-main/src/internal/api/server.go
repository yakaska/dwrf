package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"tinyurl/internal/api/handlers"
	"tinyurl/internal/shorten"
)

type Server struct {
	e       *echo.Echo
	dwarfer *shorten.Service
	baseUrl string
}

func New(dwarfer *shorten.Service, baseUrl string) *Server {

	server := &Server{dwarfer: dwarfer, baseUrl: baseUrl}
	server.setupRouter()

	return server
}

func (s *Server) setupRouter() {
	s.e = echo.New()
	s.e.HideBanner = true
	s.e.Validator = NewValidator()

	s.e.Pre(middleware.RemoveTrailingSlash())
	s.e.Use(middleware.RequestID())

	s.e.POST("/dwarfy", handlers.HandleDwarfy(s.dwarfer, s.baseUrl))
	s.e.GET("/:id", handlers.HandleRedirect(s.dwarfer))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.e.ServeHTTP(w, r)
}
