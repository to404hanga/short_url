package services

import (
	"context"
	"short_url/repository"
)

type InService struct {
	repo *repository.InRepository
}

func NewInService(repo *repository.InRepository) *InService {
	return &InService{
		repo: repo,
	}
}

func (s *InService) Create(ctx context.Context, originUrl string) error {

}
