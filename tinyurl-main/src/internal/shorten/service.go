package shorten

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"time"
	"tinyurl/internal/model"
)

type Storage interface {
	Save(ctx context.Context, link model.Link) (*model.Link, error)
	Load(ctx context.Context, linkId string) (*model.Link, error)
	AddVisits(ctx context.Context, linkId string) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Shorten(ctx context.Context, linkInput model.LinkInput) (*model.Link, error) {
	linkId := uuid.New().ID()
	var shortLink string
	if linkInput.Id != "" {
		shortLink = linkInput.Id
	} else {
		shortLink = Shorten(linkId)
	}

	dbLink := model.Link{
		Short:     shortLink,
		Long:      linkInput.URL,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	savedLink, err := s.storage.Save(ctx, dbLink)
	if err != nil {
		return nil, err
	}
	return savedLink, nil
}

func (s *Service) Get(ctx context.Context, id string) (*model.Link, error) {
	shortLink, err := s.storage.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	return shortLink, nil
}

func (s *Service) Redirect(ctx context.Context, id string) (string, error) {
	shortLink, err := s.storage.Load(ctx, id)

	if err != nil {
		return "", err
	}

	if err := s.storage.AddVisits(ctx, id); err != nil {
		slog.Error(fmt.Sprintf("failed to increment visits for identifier %q: %v", id, err))
	}

	return shortLink.Long, nil
}
