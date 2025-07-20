package v1

import (
	"log"
	"net/http"

	"github.com/azoma13/marketplace-service/internal/service"
	"github.com/labstack/echo/v4"
)

type signInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

type authRoutes struct {
	authService service.Auth
}

func newAuthRoutes(g *echo.Group, authService service.Auth) {
	r := authRoutes{
		authService: authService,
	}

	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
}

func (r *authRoutes) signUp(c echo.Context) error {
	var input signInput

	if err := c.Bind(&input); err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.authService.CreateUser(c.Request().Context(), service.AuthCreateUserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	}

	return c.JSON(http.StatusCreated, response{
		Id:       id,
		Username: input.Username,
	})
}

func (r *authRoutes) signIn(c echo.Context) error {
	var input signInput

	if err := c.Bind(&input); err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), service.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid username or password")
		return err
	}

	type response struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusOK, response{
		Token: token,
	})
}
