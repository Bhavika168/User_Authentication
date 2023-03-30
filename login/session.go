package login

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

func SetSessionKey(name string) string {
	ctx := context.Background()
	sessionId := uuid.New().String()

	sessionKey := "session:" + sessionId
	sessiondata := name

	err := rdb.Set(ctx, sessionKey, sessiondata, time.Hour).Err()
	if err != nil {
		panic(err)
	}
	return sessionKey
}

func GetSessionKey(sessionkey string) string {

	ctx := context.Background()
	name, err := rdb.Get(ctx, sessionkey).Result()
	if err != nil {
		panic(err)
	}
	return name
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
