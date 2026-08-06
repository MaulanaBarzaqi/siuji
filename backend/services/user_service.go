package services

import (
	"errors"
	"mime/multipart"
	"siuji-backend/models"
	"siuji-backend/repositories"
	"siuji-backend/utils"

	"github.com/xuri/excelize/v2"
)

type UserService interface {
	CreateParticipant(req UserRequest) (*models.User, error)
	ImportParticipantsFromExcel(file multipart.File) (int, error)
	UpdateParticipant(userID uint, req UserRequest) (*models.User, error)
	GetAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
	DeleteParticipant(userID uint) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

type UserRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	NIM        string `json:"nim"`
	University string `json:"university"`
}

func (s *userService) CreateParticipant(req UserRequest) (*models.User, error) {
	if s.userRepo.EmailExists(req.Email) {
		return nil, errors.New("email already registered")
	}
	if req.NIM == "" {
		return nil, errors.New("nim cannot be empty")
	}
	hashedPassword, err := utils.HashPassword(req.NIM)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user := models.User{
		Name: req.Name,
		Email: req.Email,
		Password: hashedPassword,
		Role: "participant",
		NIM: req.NIM,
		University: req.University,
	}
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) ImportParticipantsFromExcel(file multipart.File) (int, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return 0, errors.New("failed to read excel file")
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return 0, errors.New("failed to read rows from excel")
	}
	if len(rows) <= 1 {
		return 0, errors.New("excel file is empty or only contains header")
	}

	var usersToCreate []models.User
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 3 {
			continue
		}
		name := row[0]
		email := row[1]
		nim := row[2]
		univ := ""
		if len(row) > 3 {
			univ = row[3]
		}
		if email == "" || nim == "" || s.userRepo.EmailExists(email) {
			continue
		}
		hashedPassword, err := utils.HashPassword(nim)
		if err != nil {
			continue
		}
		user := models.User{
			Name: name,
			Email: email,
			Password: hashedPassword,
			Role: "participant",
			NIM: nim,
			University: univ,
		}
		usersToCreate = append(usersToCreate, user)
	}
	if len(usersToCreate) == 0 {
		return 0, errors.New("no valid users found to import (check for duplicate emails or missing NIM)")
	}

	// Simpan secara bulk ke database
	err = s.userRepo.CreateBulk(usersToCreate)
	if err != nil {
		return 0, err
	}

	return len(usersToCreate), nil
}

func (s *userService) UpdateParticipant(userID uint, req UserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if s.userRepo.EmailExistsExceptUser(req.Email, userID) {
		return nil, errors.New("email already used by another user")
	}

	user.Name = req.Name
	user.Email = req.Email
	user.University = req.University

	if user.NIM != req.NIM {
		user.NIM = req.NIM
		hashedPassword, err := utils.HashPassword(req.NIM)
		if err != nil {
			return nil, errors.New("failed to hash new password from nim")
		}
		user.Password = hashedPassword
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	return s.userRepo.GetAllPagination(filter, sort, limit, offset)
}

func (s *userService) DeleteParticipant(userID uint) error {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.Delete(userID)
}