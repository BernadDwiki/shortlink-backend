package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/BernadDwiki/shortlink-backend/internal/dto"
	"github.com/BernadDwiki/shortlink-backend/internal/model"
	"github.com/BernadDwiki/shortlink-backend/internal/repository"
)

var (
	ErrInvalidSlug      = errors.New("slug must be 3-50 characters, alphanumeric and hyphens only")
	ErrReservedSlug     = errors.New("slug is a reserved word")
	ErrSlugAlreadyTaken = errors.New("slug already taken")
	ErrLinkNotFound     = errors.New("link not found")
	ErrUnauthorized     = errors.New("unauthorized")
)

var slugRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

var reservedSlugs = map[string]bool{
	"api":       true,
	"login":     true,
	"register":  true,
	"dashboard": true,
}

type LinkService struct {
	repo    *repository.LinkRepository
	baseURL string
}

func NewLinkService(repo *repository.LinkRepository, baseURL string) *LinkService {
	return &LinkService{repo: repo, baseURL: baseURL}
}

func (s *LinkService) generateSlug() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *LinkService) validateSlug(slug string) error {
	if len(slug) < 3 || len(slug) > 50 {
		return ErrInvalidSlug
	}

	// Check if slug is alphanumeric and hyphens only
	if !slugRegex.MatchString(slug) {
		return ErrInvalidSlug
	}

	// Check if slug is reserved
	if reservedSlugs[strings.ToLower(slug)] {
		return ErrReservedSlug
	}

	return nil
}

func (s *LinkService) CreateLink(ctx context.Context, userID int, req dto.CreateLinkRequest) (*dto.LinkResponse, error) {
	var slug string

	// If custom slug provided, validate it
	if req.Slug != "" {
		if err := s.validateSlug(req.Slug); err != nil {
			return nil, err
		}
		slug = req.Slug
	} else {
		// Generate random slug, ensure uniqueness
		for i := 0; i < 10; i++ {
			slug = s.generateSlug()
			_, err := s.repo.GetLinkBySlug(ctx, slug)
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
	}

	link, err := s.repo.CreateLink(ctx, userID, req.OriginalURL, slug)
	if err != nil {
		if strings.Contains(err.Error(), "slug already taken") {
			return nil, ErrSlugAlreadyTaken
		}
		return nil, err
	}

	return &dto.LinkResponse{
		ID:       link.ID,
		Original: link.OriginalURL,
		Slug:     link.Slug,
		ShortURL: fmt.Sprintf("%s/%s", s.baseURL, link.Slug),
	}, nil
}

func (s *LinkService) GetUserLinks(ctx context.Context, userID int) ([]dto.LinkResponse, error) {
	links, err := s.repo.GetLinksByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.LinkResponse
	for _, link := range links {
		responses = append(responses, dto.LinkResponse{
			ID:       link.ID,
			Original: link.OriginalURL,
			Slug:     link.Slug,
			ShortURL: fmt.Sprintf("%s/%s", s.baseURL, link.Slug),
		})
	}

	return responses, nil
}

func (s *LinkService) DeleteLink(ctx context.Context, linkID, userID int) error {
	return s.repo.DeleteLink(ctx, linkID, userID)
}

func (s *LinkService) GetLinkBySlug(ctx context.Context, slug string) (*model.Link, error) {
	return s.repo.GetLinkBySlug(ctx, slug)
}
