package repositories

import (
	"errors"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id uint) (*models.User, error)
	FindByPublicID(publicID string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	UpdatePassword(userID uint, hashedPassword string) error
	FindAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
	Delete(publicID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
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
	err := r.db.Where("public_id = ?", publicID).First(&user).Error
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
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).Update("password", hashedPassword).Error
}

func (r *userRepository) FindAllPagination(filter, sort string, limit, offset int) ([]models.User,int64, error)  {
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
	db := r.db.Model(&models.User{})
	if filter != "" {
		filterPattern := "%" + filter + "%"
		// Diperluas agar bisa memfilter berdasarkan nama, email, nim, atau kampus
		db = db.Where("name ILIKE ? OR email ILIKE ? OR nim ILIKE ? OR university ILIKE ?", filterPattern, filterPattern, filterPattern, filterPattern)
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

func (r *userRepository) Delete(publicID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Where("public_id = ?", publicID).First(&user).Error; err != nil {
			return err
		}
		// Hapus juga relasi participant_periods milik user ini agar tidak menjadi orphaned data
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.ParticipantPeriod{}).Error; err != nil {
			return err
		}
		// Hapus user (GORM otomatis menangani soft delete jika field DeletedAt ada pada model User)
		return tx.Delete(&user).Error
	})
}
            
