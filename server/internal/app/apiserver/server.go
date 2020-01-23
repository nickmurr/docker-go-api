package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nickmurr/go-http-rest-api/model"
	"github.com/nickmurr/go-http-rest-api/store"
	"strings"

	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	sessionName        = "go-docker-api"
	ctxKeyUser  ctxKey = iota
)

var (
	errIncorrentCredentials = errors.New("Incorrent credentials")
	errUnathorized          = errors.New("Unathorized")
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
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		w.Header().Add("Content-Type", "application/json")
		token = w.Header().Get("Authorization")
		// tokens, ok := r.Header["Authorization"]
		if  len(token) >= 1 {
			token = strings.TrimPrefix(token, "Bearer ")
		}
		fmt.Println("token:", token)
		if token == "" {
			s.error(w, r, http.StatusUnauthorized, errUnathorized)
			return
		}

		userId, err := model.CheckJwtToken(token)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnathorized)
		}

		u, err := s.store.User().FindById(userId)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errUnathorized)
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

		// session, err := s.sessionStore.Get(r, sessionName)
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// }
		//
		// session.Values["token"] = token
		// err = s.sessionStore.Save(r, w, session)
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// }

		s.respond(w, r, http.StatusOK, token)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
