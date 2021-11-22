package service

import (
	"MyLog-M/errors"
	"MyLog-M/internal/domain"
	"MyLog-M/internal/repository"
)

//go:generate mockgen -source=./mylog.go -destination=./mock/mock.go -package=mock
type service interface {
	//Tail(limit int64) (*[]domain.Data, error)
	Get(id int64) (*domain.Data, error)
	Insert(data domain.Data) (int64, error)
	Tail(limit int64) (*[]domain.Data, error)
}

// Service/Usecase/Controller
type Service struct {
	repo repository.LogRepo
}

func New(repo repository.LogRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(id int64) (*domain.Data, error) {
	res, err := s.repo.Get(id)
	if err != nil {
		return &domain.Data{}, errors.Error{
			ErrorCode:     "DATABASE_ERROR",
			InternalError: err,
			Message:       "error occurred in database",
		}
	}
	return res, nil
}

func (s *Service) Insert(data domain.Data) (int64, error) {
	res, err := s.repo.Insert(&data)
	if err != nil {
		return res, errors.Error{
			ErrorCode:     "DATABASE_ERROR",
			InternalError: err,
			Message:       "error occurred in database",
		}
	}
	return res, nil
}

func (s *Service) Tail(limit int64) (*[]domain.Data, error) {
	res, err := s.repo.Tail(limit)
	if err != nil {
		return &[]domain.Data{}, errors.Error{
			ErrorCode:     "DATABASE_ERROR",
			InternalError: err,
			Message:       "error occurred in database",
		}
	}
	return res, nil
}
