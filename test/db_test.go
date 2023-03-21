package db

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Secret   string `json:"secret"`
}

func TestDb(t *testing.T) {
	dsn := "host=172.27.192.1 user=postgres password=123456 dbname=first_db port=9000 sslmode=disable"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	if err := db.WithContext(ctx).AutoMigrate(&User{}); err != nil {
		panic("failed to migrate database schema")
	}

	password := "mypassword"
	hpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	const charset = "abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	secretcode := make([]byte, 20)
	for i := range secretcode {
		secretcode[i] = charset[seededRand.Intn(len(charset))]
	}

	user := User{Username: "Divya", Password: string(hpassword), Secret: string(secretcode)}
	if err := db.WithContext(ctx).Create(&user).Error; err != nil {
		panic("failed to create user")
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

}
