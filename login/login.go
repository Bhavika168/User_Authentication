package login

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db.AutoMigrate(&User{})
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if CheckUserPassword(u.Username, u.Password) {
		secret := SetSecretKey(u.Username)
		message := Message{Status: "Successful", YourKey: secret, Message: "Enter Your OTP"}
		jsonStr, _ := json.Marshal(message)
		w.Write(jsonStr)

	} else {
		w.WriteHeader(http.StatusNotFound)
		data := Data{Status: "Unsuccessful", Message: "User not found."}
		jsonStr, _ := json.Marshal(data)
		w.Write(jsonStr)
		return
	}
}

func CheckUserPassword(Username, Password string) bool {
	db.AutoMigrate(&User{})

	var user User
	if err := db.Where("username = ?", Username).First(&user).Error; err != nil {
		return false
	} else {
		storedHash := user.Password
		plaintextPassword := Password
		err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(plaintextPassword))
		if err != nil {
			return false
		}
	}
	return true
}
