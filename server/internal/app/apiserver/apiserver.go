package apiserver

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start
func (s *APIServer) Start() error {
	err := s.configureLogger()
	if err != nil {
		return err
	}
	s.configureRouter()
	s.logger.Info("starting api server")
	fmt.Println(os.Getenv("BIND_ADDR"))
	return http.ListenAndServe(os.Getenv("BIND_ADDR"), s.router)
}

func (s APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *APIServer) handleHello() http.HandlerFunc {
	// ...

	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Hello")
	}

}
