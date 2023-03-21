package login

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var rdb *redis.Client

type User struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Secret   string `json:"secret"`
}

type Data struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Message struct {
	Status  string `json:"status"`
	YourKey string `json:"yourkey"`
	Message string `json:"message"`
}

func InitialiseDb() *gorm.DB {
	dsn := "host=172.27.192.1 user=postgres password=123456 dbname=first_db port=9000 sslmode=disable"
	fmt.Println(dsn)
	db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db
}

func InitialiseRedis() *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}
