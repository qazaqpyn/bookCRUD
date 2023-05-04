package service

import "github.com/qazaqpyn/bookCRUD/pkg/repository"

type Authorization interface {
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
