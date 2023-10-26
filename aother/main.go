package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const portNumber = ":8081"

type getUserResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	RoleID    int       `json:"role_id"`
}

func main() {
	response := getUserResponse{
		ID:        uuid.New(),
		FirstName: "ali",
		LastName:  "aristanov",
		Email:     "eine@gmail.com",
		Username:  "nanmenkaimak",
		Password:  hashPassword("Nanmenkaimak1*"),
		RoleID:    1,
	}

	//url := fmt.Sprintf("localhost:8081/api/user/%s", response.Username)

	http.HandleFunc("/api/user/nanmenkaimak", func(w http.ResponseWriter, r *http.Request) {
		res, _ := json.Marshal(response)
		w.Write(res)
	})

	addr := fmt.Sprintf(":%d", 8081)
	fmt.Printf("Server is listening on port %d...\n", 8081)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}
