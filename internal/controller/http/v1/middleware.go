package v1

import (
	"log"
	"net/http"
	"strings"

	"github.com/azoma13/marketplace-service/internal/service"
	"github.com/labstack/echo/v4"
)

const (
	userIdCtx = "userId"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (h *AuthMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := getToken(c.Request())
		if !ok {
			if c.Path() == "/api/v1/advertise/feed-ad" {
				return next(c)
			}
			log.Println("AuthMiddleware.UserIdentity: bearerToken: invalid auth header")
			newErrorResponse(c, http.StatusUnauthorized, "authorization error")
			return nil
		}

		userId, err := h.authService.ParseToken(token)
		if err != nil {
			log.Printf("AuthMiddleware.UserIdentity: h.authService.ParseToken: cannot parse token: %v", err)
			newErrorResponse(c, http.StatusUnauthorized, "authorization error")
			return err
		}

		c.Set(userIdCtx, userId)

		return next(c)
	}
}

func getToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
