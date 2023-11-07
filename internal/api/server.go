package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	ctxKeyUserID
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

	private.Post("/work", s.createWork)
	private.Put("work/:id", s.updateWorkData)
	private.Delete("work/:id", s.deleteWork)

	private.Get("/works", s.getWorks)

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

func (s *server) createWork(c *fiber.Ctx) error {
	type request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	req := &request{}
	if err := c.BodyParser(req); err != nil {
		s.error(c, http.StatusBadRequest, err)
		return err
	}

	w := &model.Work{
		Name:        req.Name,
		Description: req.Description,
		User_ID:     c.UserContext().Value(ctxKeyUserID).(int),
	}

	if err := s.store.Work().Create(w); err != nil {
		s.error(c, http.StatusInternalServerError, err)
		return err
	}

	return s.respond(c, http.StatusCreated, w)
}

func (s *server) getWorks(c *fiber.Ctx) error {
	w, err := s.store.Work().Get(c.UserContext().Value(ctxKeyUserID).(int))
	if err != nil {
		s.error(c, http.StatusInternalServerError, err)
		return err
	}
	return s.respond(c, http.StatusOK, w)
}

func (s *server) updateWorkData(c *fiber.Ctx) error {
	type request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	req := &request{}
	if err := c.BodyParser(req); err != nil {
		s.error(c, http.StatusBadRequest, err)
		return err
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.error(c, http.StatusBadRequest, err)
		return err
	}

	w := &model.Work{
		Name:        req.Name,
		Description: req.Description,
		User_ID:     c.UserContext().Value(ctxKeyUserID).(int),
	}

	if err := s.store.Work().Update(w, id); err != nil {
		s.error(c, http.StatusInternalServerError, err)
		return err
	}

	return s.respond(c, http.StatusOK, w)
}

func (s *server) deleteWork(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.error(c, http.StatusBadRequest, err)
		return err
	}
	user_id := c.UserContext().Value(ctxKeyUserID).(int)
	if err := s.store.Work().Delete(id, user_id); err != nil {
		s.error(c, http.StatusInternalServerError, err)
		return err
	}
	return s.respond(c, http.StatusOK, nil)
}

func (s *server) error(c *fiber.Ctx, code int, err error) {
	msg := fmt.Sprintf("STATUS CODE: %d with error %s", code, err)
	s.logger.Error(msg)
	_ = s.respond(c, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(c *fiber.Ctx, code int, data interface{}) error {
	c.Response().Header.SetStatusCode(code)
	if data != nil {
		_ = json.NewEncoder(c.Response().BodyWriter()).Encode(data)
	}
	return nil
}
