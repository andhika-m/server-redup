package middleware

import (
	"context"
	"net/http"
	"redup/internal/model/constant"
	"redup/internal/usecase/redup"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type authMiddleware struct {
	redupUsecase redup.Usecase
}

func LoadMiddlewares(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{

		AllowOrigins: []string{"*"},
	}))
}

func GetAuthMiddleware(redupUsecase redup.Usecase) *authMiddleware {
	return &authMiddleware{
		redupUsecase: redupUsecase,
	}
}

func (am *authMiddleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionData, err := GetSessionData(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     401,
				Message:  err.Error(),
				Internal: err,
			}
		}

		userID, err := am.redupUsecase.CheckSession(sessionData)
		if err != nil {
			return &echo.HTTPError{
				Code:     401,
				Message:  err.Error(),
				Internal: err,
			}
		}

		authContext := context.WithValue(c.Request().Context(), constant.AuthContextKey, userID)
		c.SetRequest(c.Request().WithContext(authContext))

		if err := next(c); err != nil {
			return err
		}

		return nil
	}
}

func (m *authMiddleware) CheckUsers(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, "missing or malformed jwt")
		}

		token := strings.Split(authHeader, " ")[1]

		userID, err := m.redupUsecase.VerifySession(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "invalid or expired jwt")
		}

		c.Set("userID", userID)
		c.Set("token", token)

		return next(c)
	}
}
