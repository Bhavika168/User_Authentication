package login

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	db.AutoMigrate(&User{})
	var u User
	json.NewDecoder(r.Body).Decode(&u)

	hpassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	secretcode := GenerateSecretUserKey()

	user := User{Username: u.Username, Password: string(hpassword), Secret: secretcode}
	if err := db.WithContext(ctx).Create(&user).Error; err != nil {
		panic("failed to create user")
	}

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
}

func GenerateSecretUserKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	secretcode := make([]byte, 20)
	for i := range secretcode {
		secretcode[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(secretcode)
}

func GenerateQRcode(username string) string {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      username,
		AccountName: username + "_qr",
	})

	img, err := key.Image(200, 200)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding QR code: %v\n", err)
		os.Exit(1)
	}

	file, _ := os.Create("qr.png")
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding QR code image: %v\n", err)
		os.Exit(1)
	}
	return key.Secret()
}
