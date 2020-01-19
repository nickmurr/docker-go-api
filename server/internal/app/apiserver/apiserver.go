package apiserver

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nickmurr/go-http-rest-api/store"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
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

	_ = s.configureStore()

	s.logger.Info("starting api server")

	fmt.Println(os.Getenv("BIND_ADDR"))

	return http.ListenAndServe(os.Getenv("BIND_ADDR"), s.router)
}

func (s *APIServer) configureLogger() error {
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

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.logger.Info("Opened db successfuly")

	s.store = st
	return nil
}

func (s *APIServer) handleHello() http.HandlerFunc {
	// ...

	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Hello")
	}

}
