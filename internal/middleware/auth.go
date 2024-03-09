package middleware

import (
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	// authUC authUsecase.UsecaseI
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().URL.Path == "/signup" || c.Request().URL.Path == "/signin" ||
			c.Request().URL.Path == "/auth" || c.Request().URL.Path == "/prometheus" ||
			c.Request().URL.Path == "/favicon.ico" {
			return next(c)
		}

		// TODO логика с авторизацией

		c.Set("user_id", int64(1))
		return next(c)
	}
}
