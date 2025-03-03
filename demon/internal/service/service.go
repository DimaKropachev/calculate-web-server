package service

import (
	"github.com/DimaKropachev/calculate-web-server/demon/internal/model"
)

type Api interface {
	Get() (*model.Task, error)
	Give(resp *model.Response) error
}

type Service struct {
	api Api
}

func NewService(api Api) *Service {
	return &Service{
		api: api,
	}
}

func (s *Service) Get() (*model.Task, error) {
	return s.api.Get()
}

func (s *Service) Give(resp *model.Response) error {
	return s.api.Give(resp)
}