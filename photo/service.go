package photo

import (
	"context"

	"github.com/anggi-susanto/go-face-detection-be/domain"
	"github.com/anggi-susanto/go-face-detection-be/internal/repository/mongo"
)

type Service interface {
	Save(ctx context.Context, photo *domain.Photo) error
	CheckResult(ctx context.Context, id string) (*domain.Photo, error)
	GetPhoto(ctx context.Context, id string) (*domain.Photo, error)
}

type service struct {
	photoRepository mongo.PhotoRepository
}

func NewService(photoRepository mongo.PhotoRepository) Service {
	return &service{
		photoRepository: photoRepository,
	}
}

func (s *service) Save(ctx context.Context, photo *domain.Photo) error {
	return s.photoRepository.Create(ctx, photo)
}

func (s *service) CheckResult(ctx context.Context, id string) (*domain.Photo, error) {
	return s.photoRepository.FindByID(ctx, id)
}

func (s *service) GetPhoto(ctx context.Context, id string) (*domain.Photo, error) {
	return s.photoRepository.FindByID(ctx, id)
}
