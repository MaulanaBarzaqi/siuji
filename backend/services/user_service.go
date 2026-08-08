package services

import (
	"siuji-backend/models"
	"siuji-backend/repositories"
)

type UserService interface {
	GetAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
	// DeleteParticipant(userID uint) error
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

// func (s *userService) DeleteParticipant(userID uint) error {
// 	_, err := s.userRepo.FindByID(userID)
// 	if err != nil {
// 		return errors.New("user not found")
// 	}
// 	return s.userRepo.Delete(userID)
// }