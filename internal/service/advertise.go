package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/azoma13/marketplace-service/config"
	"github.com/azoma13/marketplace-service/internal/entity"
	"github.com/azoma13/marketplace-service/internal/repo"
	"github.com/azoma13/marketplace-service/internal/repo/pgdb"
)

type AdvertiseService struct {
	advertiseRepo repo.Advertise
	userRepo      repo.User
}

func NewAdvertiseService(advertiseRepo repo.Advertise, userRepo repo.User) *AdvertiseService {
	return &AdvertiseService{
		advertiseRepo: advertiseRepo,
		userRepo:      userRepo,
	}
}

func (s *AdvertiseService) CreateAdvertise(ctx context.Context, input AdvertiseCreateNewAdvertiseInput) (entity.Advertise, error) {
	err := ValidImage(input.Image)
	if err != nil {
		return entity.Advertise{}, fmt.Errorf("AdvertiseService.CreateAdvertise - ValidImage: %v", err)
	}

	advertise := &entity.Advertise{
		Title:       input.Title,
		Description: input.Description,
		Image:       input.Image,
		Price:       input.Price,
		UserId:      input.UserId,
	}

	adv, err := s.advertiseRepo.CreateAdvertise(ctx, *advertise)
	if err != nil {
		return entity.Advertise{}, fmt.Errorf("AdvertiseService.CreateAdvertise - s.advertiseRepo.CreateAdvertise: %v", err)
	}

	return adv, nil
}

func (s *AdvertiseService) GetFeedAdvertise(ctx context.Context, input AdvertiseGetFeedAdInput) ([]*entity.FeedAdvertises, error) {
	params, err := ParseParam(input)
	if err != nil {
		return nil, fmt.Errorf("AdvertiseService.CreateAdvertise - ParseParam: %v", err)
	}

	adv, err := s.advertiseRepo.GetFeedAdvertise(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("AdvertiseService.CreateAdvertise - s.advertiseRepo.CreateAdvertise: %v", err)
	}

	return adv, nil
}

func ValidImage(imageUrl string) error {
	decodedUrlFile, err := url.QueryUnescape(imageUrl)
	if err != nil {
		return err
	}
	urlFile := strings.ReplaceAll(decodedUrlFile, " ", "_")

	fileName := path.Base(urlFile)
	if idx := strings.Index(fileName, "?"); idx != -1 {
		fileName = fileName[:idx]
	}

	ext := path.Ext(fileName)
	ok := slices.Contains(config.Cfg.AllowedFileExtensions, ext)
	if !ok {
		return fmt.Errorf("extension file not allowed")
	}

	response, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("file was not found link: %s", imageUrl)
	}

	size := config.Cfg.APP.MaxImageSize * 1024 * 1024
	limitReader := io.LimitReader(response.Body, size+1)
	buffer := new(bytes.Buffer)

	n, err := buffer.ReadFrom(limitReader)
	if err != nil && err != io.EOF {
		return err
	}

	log.Println(n, size)
	if n > size {
		return fmt.Errorf("max image size has been exceeded")
	}

	return nil
}

func ParseParam(input AdvertiseGetFeedAdInput) (pgdb.AdvertiseGetFeedAdRepo, error) {
	pageInt, err := strconv.Atoi(input.Page)
	if err != nil {
		return pgdb.AdvertiseGetFeedAdRepo{}, fmt.Errorf("error convert page: %v", err)
	}
	perPageInt, err := strconv.Atoi(input.PerPage)
	if err != nil {
		return pgdb.AdvertiseGetFeedAdRepo{}, fmt.Errorf("error convert perPage: %v", err)
	}

	sorts := strings.Split(input.Sort, "_")
	var order string
	if len(sorts) == 1 {
		order = "ASC"
	} else {
		order = "DESC"
	}

	var sort string
	if sorts[0] == "price" {
		sort = "a.price"
	} else {
		sort = "a.created_at"
	}

	decodedStr, err := url.QueryUnescape(input.CurrentPrice)
	if err != nil {
		return pgdb.AdvertiseGetFeedAdRepo{}, fmt.Errorf("failed to decode price range: %w", err)
	}

	prices := strings.Split(decodedStr, ";")
	if len(prices) != 2 {
		return pgdb.AdvertiseGetFeedAdRepo{}, fmt.Errorf("invalid price format")
	}

	minPrice, err := strconv.ParseFloat(prices[0], 64)
	if err != nil {
		return pgdb.AdvertiseGetFeedAdRepo{}, fmt.Errorf("failed to parse min price: %w", err)
	}

	maxPrice, err := strconv.ParseFloat(prices[1], 64)
	if err != nil {
		return pgdb.AdvertiseGetFeedAdRepo{}, fmt.Errorf("failed to parse max price: %w", err)
	}

	return pgdb.AdvertiseGetFeedAdRepo{
		UserId:    input.UserId,
		Sort:      sort,
		SortOrder: order,
		Page:      pageInt,
		PerPage:   perPageInt,
		MinPrice:  minPrice,
		MaxPrice:  maxPrice,
	}, nil
}
