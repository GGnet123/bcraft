package server

import (
	"bcraft/docs"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

var defaultStopTimeout = time.Second * 30

type Server struct {
	httpServer *http.Server
}

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func Init() {
	docs.SwaggerInfo.Title = "Bcraft API Server"
	docs.SwaggerInfo.Description = "This is bcraft test server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	srv := new(Server)
	handler := new(Router)
	if err := srv.Run(os.Getenv("SERVER_PORT"), handler.Init()); err != nil {
		logrus.Fatal("couldn't start the server. Error: %s", err.Error())
	}
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    defaultStopTimeout,
		WriteTimeout:   defaultStopTimeout,
	}

	return s.httpServer.ListenAndServe()
}
