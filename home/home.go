package home

import (
	"Project12/login"
	"encoding/json"
	"net/http"
)

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		defer func() {
			if err := recover(); err != nil {

				w.WriteHeader(http.StatusUnauthorized)
				data := Message{Status: "Unsuccessful", Message: "Wrong Header."}
				jsonStr, _ := json.Marshal(data)
				w.Write(jsonStr)
			}
		}()

		username := login.GetSessionKey(r.Header.Get("X-Session"))
		if username != "" {
			// fmt.Fprintf(w, "User is inside home.\n")
		} else {
			panic("Wrong Header")
		}

		next.ServeHTTP(w, r)
	})
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := Message{Status: "Successful", Message: "Inside Home."}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Write(jsonStr)
}
