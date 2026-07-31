package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB		*gorm.DB
	AppConfig *Config
)

type Config struct {
	Port                string
	APPURL              string
	CORSOrigin          string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	JWTSecret           string
	JWTExpiresIn        string
	RefreshTokenExpires string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("warning: No .Env file found, using environment variables.")
	}
	
AppConfig = &Config{
		Port:                getEnv("PORT", "8080"),
		APPURL:              getEnv("APP_URL", "http://localhost:8080"),
		CORSOrigin:          getEnv("CORS_ORIGIN", "http://localhost:5173"),
		DBHost:              getEnv("DB_HOST", "db"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "your_postgres_user"),
		DBPassword:          getEnv("DB_PASSWORD", "your_postgres_password"),
		DBName:              getEnv("DB_NAME", "your_database_name"),
		JWTSecret:           getEnv("JWT_SECRET", "your_jwt_secret_here"),
		JWTExpiresIn:        getEnv("JWT_EXPIRES_IN", "2h"),
		RefreshTokenExpires: getEnv("REFRESH_TOKEN_EXPIRES", "24h"),
	}
	log.Println("environment variables loaded successfully.")
}

func getEnv(key string, fallback string) string {
	value, exist := os.LookupEnv(key)
	if exist {
		return value
	}else {
		return fallback
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func ConnectDB() {
	cfg := AppConfig

	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database instance:", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("successfully connected to PostgreSQL!")
}