package login

import (
	"encoding/base64"
	"encoding/json"
	"image"
	"image/png"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	db.AutoMigrate(&User{})
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if CheckUserPassword(u.Username, u.Password) {
		code := GenerateQRcode(u.Username)
		SetSessionQR(u.Username, code)

		file, _ := os.Open("qr.png")
		defer file.Close()
		img, _, _ := image.Decode(file)

		var buf []byte
		png.Encode(w, img)
		base64Str := base64.StdEncoding.EncodeToString(buf)
		QRcode, _ := json.Marshal(map[string]string{"qr_base64": base64Str})

		w.Header().Set("Content-Type", "application/json")
		w.Write(QRcode)

	} else {
		w.Header().Set("Content-Type", "application/json")
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
