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
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		helpers.EnableCors(w)
		helpers.LogAction("Welcome")
	})

	http.HandleFunc("/auth/login", handlers.Login(db))
	http.HandleFunc("/auth/signup", handlers.SignUp(db))
	http.HandleFunc("/auth/send-reset-password-email", handlers.SendResetMail(db))
	http.HandleFunc("/auth/reset-password", handlers.ResetPassword(db))
	http.HandleFunc("/users", handlers.SearchUsers(db))
	http.HandleFunc("/projects", handlers.AddProject(db))
	http.HandleFunc("/project/", handlers.SingleProject(db))
	http.HandleFunc("/project/member", handlers.AddMember(db))
	http.HandleFunc("/project/members", handlers.GetMembers(db))
	http.HandleFunc("/sprints", handlers.AddSprint(db))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
