package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *server) setRequestID(c *fiber.Ctx) error {
	id := uuid.New().String()
	c.Response().Header.Set("X-Request-ID", id)
	c.SetUserContext(context.WithValue(c.Context(), ctxKeyRequestID, id))
	return c.Next()
}

func (s *server) logRequest(c *fiber.Ctx) error {
	logger := s.logger
	msg := fmt.Sprintf("remote_addr=%s request_id=%s", c.Context().RemoteAddr(), c.UserContext().Value(ctxKeyRequestID))
	logger.Info(fmt.Sprintf("started %s %s \t %s", c.Context().Method(), c.Context().RequestURI(), msg))
	start := time.Now()
	msg = fmt.Sprintf("completed with %d % s in %v \t %s",
		http.StatusOK,
		http.StatusText(http.StatusOK),
		time.Since(start),
		msg)
	logger.Info(msg)
	return c.Next()
}

func (s *server) authUser(c *fiber.Ctx) error {
	claims, err := validateToken(c)
	if err != nil {
		s.error(c, http.StatusBadRequest, err)
		return err
	}

	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		s.error(c, http.StatusInternalServerError, errors.New("token expired"))
		return errors.New("token expired")
	}
	u, err := s.store.User().FindByID(int(claims["sub"].(float64)))
	if err != nil {
		s.error(c, http.StatusInternalServerError, err)
		return err
	}
	c.Context().SetUserValue(ctxKeyUser, u)
	c.Context().SetUserValue(ctxKeyUserID, u.ID)
	return c.Next()
}

func (s *server) getUserID(c *fiber.Ctx) int {
	userID := c.UserContext().Value(ctxKeyUserID).(int)
	return userID
}
