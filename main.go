package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ralphvw/sprint-planner-v2/db"
	"github.com/ralphvw/sprint-planner-v2/handlers"
	"github.com/ralphvw/sprint-planner-v2/helpers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	db := db.InitDb()

	fmt.Print("Server started at " + port + "\n")
	http.HandleFunc("/", func(http.ResponseWriter, *http.Request) {
		helpers.LogAction("Welcome")
	})

	http.HandleFunc("/auth/login", handlers.Login(db))
	http.HandleFunc("/auth/signup", handlers.SignUp(db))
	http.HandleFunc("/auth/send-reset-password-email", handlers.SendResetMail(db))
  http.HandleFunc("/auth/reset-password", handlers.ResetPassword(db))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
