package seed

import (
	"log"
	"siuji-backend/config"
	"siuji-backend/models"
	"siuji-backend/utils"

	"github.com/google/uuid"
)

func SeedAdmin() {
	var adminCount int64
	config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount > 0 {
		log.Println("admin user already exists. skiping seed to prevent duplicates.")
		return
	}
	log.Println("no admin found. seeding default admin user")
	
	password, err := utils.HashPassword("Adminsiuji2026!")
	if err != nil {
		log.Fatal("failed to hash password!", err)
	}

	admin := models.User{
		PublicId: uuid.New(),
		Name: "admin siuji",
		Email: "adminsiuji@siuji.com",
		Password: password,
		Role: "admin",
	}

	if err := config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error; err != nil {
		log.Println("fail to seed admin:", err)
	}else {
		log.Println("admin user seeded")
	}
}