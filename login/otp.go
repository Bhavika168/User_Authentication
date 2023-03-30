package login

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pquerna/otp/totp"
)

type TOTP struct {
	Username string `json:"username"`
	YourOTP  string `json:"yourotp"`
}

func CheckOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			data := Data{Status: "Unsuccessful", Message: "Wrong OTP."}
			jsonStr, _ := json.Marshal(data)
			w.Write(jsonStr)
		}
	}()

	var t TOTP
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := GetSecretKey(t.YourOTP)
	if t.Username == name {
		key := SetSessionKey(t.Username)
		data1 := Data{Status: "Successful", Message: key}
		jsonStr, _ := json.Marshal(data1)
		w.Write(jsonStr)
	} else {
		panic("Wrong OTP.")
	}
}

func SetSecretKey(name string) string {
	ctx := context.Background()
	if err := db.WithContext(ctx).AutoMigrate(&User{}); err != nil {
		panic("failed to migrate database schema")
	}

	var user User
	err1 := db.Where("username = ?", name).First(&user).Error
	if err1 != nil {
		panic(err1)
	}
	secret := user.Secret
	secretKey, _ := totp.GenerateCode(secret, time.Now())
	secretData := name

	err := rdb.Set(ctx, secretKey, secretData, time.Minute).Err()
	if err != nil {
		panic(err)
	}
	return secret
}

func GetSecretKey(code string) string {
	ctx := context.Background()
	name, err := rdb.Get(ctx, code).Result()
	if err != nil {
		panic(err)
	}
	return name
}
