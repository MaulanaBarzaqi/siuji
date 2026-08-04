package repositories

import (
	"errors"
	"siuji-backend/config"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id uint) (*models.User, error)
	FindByPublicID(publicID string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	UpdateProfile(userID uint, name, email string) error
	UpdatePassword(userID uint, hashedPassword string) error
	GetAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
	Delete(userID uint) error
	EmailExists(email string) bool
	EmailExistsExceptUser(email string, userID uint) bool
}

type userRepository struct {}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPublicID(publicID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("public_id = ?", publicID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return config.DB.Save(user).Error
}

func (r *userRepository) UpdateProfile(userID uint, name, email string) error {
	return config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(
		map[string]interface{}{
			"name": name,
			"email": email,
	}).Error
}

func (r *userRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return config.DB.Model(&models.User{}).
		Where("id = ?", userID).Update("password", hashedPassword).Error
}

func (r *userRepository) GetAllPagination(filter, sort string, limit, offset int) ([]models.User,int64, error)  {
	if limit <= 0 {
		limit = 10 
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	var users []models.User
	var total int64
	db := config.DB.Model(&models.User{})
	if filter != "" {
		filterPattern := "%" + filter + "%"
		// Diperluas agar bisa memfilter berdasarkan nama, email, nim, atau kampus
		db = db.Where("name ILIKE ? OR email ILIKE ? OR nim ILIKE ? OR university_name ILIKE ?", filterPattern, filterPattern, filterPattern, filterPattern)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	allowedSortFields := map[string]string{
		"id":          "id ASC",
		"-id":         "id DESC",
		"name":        "name ASC",
		"-name":       "name DESC",
		"email":       "email ASC",
		"-email":      "email DESC",
		"nim":         "nim ASC",
		"-nim":        "nim DESC",
		"role":        "role ASC",
		"-role":       "role DESC",
		"created_at":  "created_at ASC",
		"-created_at": "created_at DESC",
	}
	if sort == "" {
		sort = "-created_at"
	}
	if sortClause, ok := allowedSortFields[sort]; ok {
		db = db.Order(sortClause)
	} else {
		db = db.Order("created_at DESC")
	}
	err := db.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, err
}

func (r *userRepository) Delete(userID uint) error {
	return config.DB.Delete(&models.User{}, userID).Error
}
            
func (r *userRepository) EmailExists(email string) bool {
	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *userRepository) EmailExistsExceptUser(email string, userID uint) bool {
	var count int64
	config.DB.Model(&models.User{}).
		Where("email = ? AND id != ?", email, userID).
		Count(&count)
	return count > 0
}