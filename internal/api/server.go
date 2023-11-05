package api

import (
	"encoding/json"
	"net/http"

	"github.com/artemiyKew/todo-list-rest-api/internal/model"
	"github.com/artemiyKew/todo-list-rest-api/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

type key int

const (
	ctxKeyRequestID key = iota
	ctxKeyUser
)

type server struct {
	router *fiber.App
	logger *zap.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: fiber.New(),
		logger: zap.NewNop(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)

	s.router.Post("/registration", s.createUser)
	s.router.Post("/login", s.loginUser)

	private := s.router.Group("/auth")
	private.Use(s.authUser)
	private.Get("/whoami", s.whoami)
}

func (s *server) createUser(c *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := &request{}

	if err := c.BodyParser(req); err != nil {
		s.error(c, http.StatusBadRequest, err)
		return err
	}
	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := s.store.User().Create(u); err != nil {
		s.error(c, http.StatusUnprocessableEntity, err)
		return err
	}
	u.Sanitize()

	return s.respond(c, http.StatusCreated, u)
}

func (s *server) loginUser(c *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := &request{}
	if err := c.BodyParser(&req); err != nil {
		s.error(c, http.StatusBadRequest, err)
		return err
	}

	u, err := s.store.User().FindByEmail(req.Email)
	if err != nil || !u.ComparePasswords(req.Password) {
		s.error(c, http.StatusInternalServerError, err)
		return err
	}

	token, err := generateJWT(u)
	if err != nil {
		s.error(c, http.StatusInternalServerError, err)
		return err
	}

	c.Response().Header.Set("Token", token)
	return s.respond(c, http.StatusAccepted, u)
}

func (s *server) whoami(c *fiber.Ctx) error {
	return s.respond(c, http.StatusOK, c.UserContext().Value(ctxKeyUser))
}

func (s *server) error(c *fiber.Ctx, code int, err error) {
	s.respond(c, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(c *fiber.Ctx, code int, data interface{}) error {
	c.Response().Header.SetStatusCode(code)
	if data != nil {
		json.NewEncoder(c.Response().BodyWriter()).Encode(data)
	}
	return nil
}
