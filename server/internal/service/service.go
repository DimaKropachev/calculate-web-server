package service

import (
	"github.com/DimaKropachev/calculate-web-server/server/internal/models"
)

type Repository interface {
	Add(expression string) string
	Get(Id string) (*models.Expression, error)
	GetAll() []*models.Expression
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(expression string) string {
	return s.repo.Add(expression)
}

func (s *Service) Get(Id string) (*models.Expression, error) {
	return s.repo.Get(Id)
}

func (s *Service) GetAll() []*models.Expression {
	return s.repo.GetAll()
}
