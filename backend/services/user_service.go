package services

import (
	"errors"
	"siuji-backend/models"
	"siuji-backend/repositories"
)

type UserService interface {
	GetAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
	GetUserByPublicID(publicID string) (*models.User, error)
	DeleteUser(publicID string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	return s.userRepo.FindAllPagination(filter, sort, limit, offset)
}

func (s *userService) GetUserByPublicID(publicID string) (*models.User, error) {
	user, err := s.userRepo.FindByPublicID(publicID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) DeleteUser(publicID string) error {
	_, err := s.userRepo.FindByPublicID(publicID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.Delete(publicID)
}