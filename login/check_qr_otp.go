package login

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pquerna/otp/totp"
)

func CheckQrOTP(w http.ResponseWriter, r *http.Request) {
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

func SetSessionQR(username, code string) string {

	ctx := context.Background()
	secretKey, _ := totp.GenerateCode(code, time.Now())
	secretData := username

	err := rdb.Set(ctx, secretKey, secretData, time.Hour).Err()
	if err != nil {
		panic(err)
	}
	return code

}

func GetSessionQR(code string) string {

	ctx := context.Background()
	name, err := rdb.Get(ctx, code).Result()
	if err != nil {
		panic(err)
	}
	return name
}
