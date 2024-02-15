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

	mux := http.NewServeMux()

	db := db.InitDb()

	fmt.Print("Server started at " + port + "\n")
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		helpers.EnableCors(w)
		helpers.LogAction("Welcome")
	})

	mux.HandleFunc("/auth/login", handlers.Login(db))
	mux.HandleFunc("/auth/signup", handlers.SignUp(db))
	mux.HandleFunc("/auth/send-reset-password-email", handlers.SendResetMail(db))
	mux.HandleFunc("/auth/reset-password", handlers.ResetPassword(db))
	mux.HandleFunc("/users", handlers.SearchUsers(db))
	mux.HandleFunc("GET /projects", handlers.GetAllProjects(db))
	mux.HandleFunc("POST /projects", handlers.AddProject(db))
	mux.HandleFunc("GET /project/{id}", handlers.GetSingleProject(db))
	mux.HandleFunc("POST /projects/members", handlers.AddMember(db))
	mux.HandleFunc("DELETE /projects/members", handlers.DeleteMember(db))
	mux.HandleFunc("GET /projects/members", handlers.GetMembers(db))
	mux.HandleFunc("GET /sprints", handlers.GetSprints(db))
	mux.HandleFunc("POST /sprints", handlers.AddSprint(db))
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
