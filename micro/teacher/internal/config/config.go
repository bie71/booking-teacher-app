package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	MysqlDSN          string
	DbMaxConection    int
	DbMaxIdleConns    int
	DbConnMaxLifetime time.Duration
	AppPort           string
	GIN_MODE          string
	JWT               JWT
	Client            Client
	ServiceBooking    Service
	IsNFT             bool
}

type JWT struct {
	SecretKey     string
	TokenDuration int
}

type Service struct {
	Host string
	Port string
}

type Client struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	Region     string
	BucketName string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	return &Config{
		MysqlDSN:          os.Getenv("MYSQL_DSN"),
		AppPort:           os.Getenv("APP_PORT"),
		GIN_MODE:          os.Getenv("GIN_MODE"),
		DbMaxConection:    cast.ToInt(os.Getenv("DB_MAX_CONECTION")),
		DbMaxIdleConns:    cast.ToInt(os.Getenv("DB_MAX_IDLE_CONNS")),
		DbConnMaxLifetime: cast.ToDuration(os.Getenv("DB_CONN_MAX_LIFETIME")),
		JWT: JWT{
			SecretKey:     os.Getenv("JWT_SECRET_KEY"),
			TokenDuration: cast.ToInt(os.Getenv("JWT_TOKEN_DURATION")),
		},
		Client: Client{
			Endpoint:   os.Getenv("CLIENT_ENDPOINT"),
			AccessKey:  os.Getenv("CLIENT_ACCESS_KEY"),
			SecretKey:  os.Getenv("CLIENT_SECRET_KEY"),
			Region:     os.Getenv("CLIENT_REGION"),
			BucketName: os.Getenv("CLIENT_BUCKET_NAME"),
		},
		ServiceBooking: Service{
			Host: os.Getenv("SERVICE_BOOKING_HOST"),
		},
		IsNFT: cast.ToBool(os.Getenv("IS_NFT")),
	}
}

func InitDB(c *Config) *gorm.DB {
	dsn := c.MysqlDSN

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Set max connection
	sqlDB, err := database.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(c.DbMaxConection)                     // Maksimal 20 koneksi terbuka ke DB
	sqlDB.SetMaxIdleConns(c.DbMaxIdleConns)                     // Maksimal 10 koneksi idle
	sqlDB.SetConnMaxLifetime(c.DbConnMaxLifetime * time.Minute) // Maksimal lifetime koneksi 30 menit

	return database
}
