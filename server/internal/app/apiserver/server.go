package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nickmurr/go-http-rest-api/model"
	"github.com/nickmurr/go-http-rest-api/store"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	sessionName        = "go-docker-api"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrentCredentials = errors.New("Incorrent credentials")
	errUnathorized          = errors.New("Unathorized")
	errWrongTokenFormat     = errors.New("Invalid/Malformed auth token")
	mySigningKey            = []byte("secret-go-api")
	customJwtMiddleware     = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
)

type ctxKey int8

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	// Private
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
}

func (s *server) handleWhoami() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	})
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})

		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s, in %v seconds",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start).Seconds())
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		w.Header().Add("Content-Type", "application/json")

		token = r.Header.Get("Authorization")

		tokens := strings.Split(token, " ")
		if len(tokens) != 2 {
			s.error(w, r, http.StatusForbidden, errWrongTokenFormat)
			return
		}

		if len(tokens) >= 1 {
			// token = strings.TrimPrefix(token, "Bearer ")
			token = tokens[1]
		}
		if token == "" {
			s.error(w, r, http.StatusUnauthorized, errUnathorized)
			return
		}

		userId, err := model.CheckJwtToken(token)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnathorized)
			return
		}

		u, err := s.store.User().FindById(userId)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnathorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)

	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrentCredentials)
			return
		}

		token, _, err := u.TokenBack(mySigningKey)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		s.respond(w, r, http.StatusOK, token)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.Header().Add("Content-Type", "application/json")
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
