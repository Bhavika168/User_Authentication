package main

import (
	"Project12/home"
	"Project12/login"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	login.InitialiseDb()
	login.InitialiseRedis()

	router := mux.NewRouter()

	signupRouter := router.PathPrefix("/signup").Subrouter()
	signupRouter.HandleFunc("", login.Signup).Methods(http.MethodGet)
	signupRouter.HandleFunc("/otp", login.CheckOTP).Methods(http.MethodGet)

	loginRouter := router.PathPrefix("/login").Subrouter()
	loginRouter.HandleFunc("", login.Login).Methods(http.MethodGet)
	loginRouter.HandleFunc("/otp", login.CheckOTP).Methods(http.MethodGet)

	homeRouter := router.PathPrefix("/home").Subrouter()
	homeRouter.Use(home.Auth)
	homeRouter.HandleFunc("", home.Home).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", router))
}
