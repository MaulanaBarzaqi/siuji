package database

import (
	"log"
	"siuji-backend/config"
	"siuji-backend/models"
)

var MigrationRegistry = []interface{} {
	&models.User{},
	&models.Period{},
	&models.Section{},
	&models.ScoreConversion{},
	&models.OTP{},
	&models.PeriodSection{},
	&models.Question{},
	&models.ParticipantPeriod{},
	&models.Option{},
	&models.SectionScore{},
	&models.AnswerKey{},
	&models.ParticipantAnswer{},
}


func Migrate() {
	log.Println("running database migration...")

	if len(MigrationRegistry) == 0 {
		log.Println("no models to migrate")
		return
	}

	err := config.DB.AutoMigrate(MigrationRegistry...)
	if err != nil {
		log.Fatal("failed to migrate databse:", err)
	}
	log.Printf("migrated %d models successfully", len(MigrationRegistry))
}