package v1

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/azoma13/marketplace-service/internal/service"
	"github.com/labstack/echo/v4"
)

type AdvertiseCreateInput struct {
	Title       string  `json:"title" validate:"required,min=4,max=255"`
	Description string  `json:"description" validate:"required,min=4,max=1024"`
	Image       string  `json:"image" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=49.99,max=999999999.99"`
}

type advertiseRoutes struct {
	advertiseService service.Advertise
}

func newAdvertiseRoutes(g *echo.Group, advertiseService service.Advertise) {
	r := advertiseRoutes{
		advertiseService: advertiseService,
	}

	g.POST("/create", r.createAdvertise)
	g.GET("/feed-ad", r.getFeedAd)
}

func (r *advertiseRoutes) createAdvertise(c echo.Context) error {
	var input AdvertiseCreateInput

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

	id, err := r.getId(c)
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return fmt.Errorf("invalid get userid")
	}

	advertise, err := r.advertiseService.CreateAdvertise(c.Request().Context(), service.AdvertiseCreateNewAdvertiseInput{
		Title:       input.Title,
		Description: input.Description,
		Image:       input.Image,
		Price:       input.Price,
		UserId:      id,
	})

	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusCreated, advertise)
}

func (r *advertiseRoutes) getId(c echo.Context) (int, error) {
	userId := c.Get(userIdCtx)
	if userId == nil {
		return 0, ErrInvalidGetUserId
	}

	id, ok := userId.(int)
	if !ok {
		return 0, fmt.Errorf("error convert userid")
	}

	return id, nil
}

func (r *advertiseRoutes) getFeedAd(c echo.Context) error {
	sort := getQueryParam(c.Request().URL.Query(), "sort", "new_desc")
	page := getQueryParam(c.Request().URL.Query(), "page", "1")
	perPage := getQueryParam(c.Request().URL.Query(), "per_page", "2")
	currentPrice := getQueryParam(c.Request().URL.Query(), "currency_price", "49.99%3B999999999.99")

	id, err := r.getId(c)
	if err != nil {
		if err != ErrInvalidGetUserId {
			log.Println(err)
			newErrorResponse(c, http.StatusBadRequest, "invalid request body")
			return fmt.Errorf("invalid get userid")
		}
	}

	feedAd, err := r.advertiseService.GetFeedAdvertise(c.Request().Context(), service.AdvertiseGetFeedAdInput{
		Sort:         sort,
		Page:         page,
		PerPage:      perPage,
		CurrentPrice: currentPrice,
		UserId:       id,
	})
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return fmt.Errorf("invalid get userid")
	}
	log.Println(feedAd)
	return c.JSON(http.StatusOK, &feedAd)
}

func getQueryParam(query url.Values, key, defaultValue string) string {
	if value := query.Get(key); value != "" {
		if ok := validatorForValueQueryParam(key, value); ok {
			return value
		}
	}

	return defaultValue
}

func validatorForValueQueryParam(key, value string) bool {
	switch key {
	case "sort":
		switch value {
		case "new":
		case "new_desc":
		case "price":
		case "price_desc":
		default:
			return false
		}
		return true
	case "page":
		_, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		return true
	case "per_page":
		switch value {
		case "10":
		case "25":
		case "50":
		default:
			return false
		}
		return true
	case "currency_price":
		return true
	default:
		return false
	}
}
