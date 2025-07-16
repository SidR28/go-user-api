package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:Honey.7511@tcp(127.0.0.1:3306)/go_api")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MySQL!")
}

func nameHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/name/")
	fmt.Fprintf(w, "Hello %s", name)
}

type Message struct {
	Message string `json:"message"`
}

func jsonNameHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/json/")
	response := Message{Message: "Hello " + name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use POST method only", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
        INSERT INTO users (firstName, lastName, email, password, age)
        VALUES (?, ?, ?, ?, ?)`,
		user.FirstName, user.LastName, user.Email, user.Password, user.Age)

	if err != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}

	fmt.Println("User created:", user.FirstName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User created Successfully.")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/get-user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	err = db.QueryRow(`
        SELECT id, firstName, lastName, email, password, age
        FROM users WHERE id = ?`, id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Age)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Use PUT method only", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/update-user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
        UPDATE users SET firstName=?, lastName=?, email=?, password=?, age=?
        WHERE id=?`,
		user.FirstName, user.LastName, user.Email, user.Password, user.Age, id)

	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	fmt.Println("User updated:", id)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User updated Successfully.")
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Use DELETE method only", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/delete-user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Database delete failed", http.StatusInternalServerError)
		return
	}

	fmt.Println("User deleted:", id)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User deleted Successfully.")
}

func main() {
	http.HandleFunc("/name/", nameHandler)
	http.HandleFunc("/json/", jsonNameHandler)
	http.HandleFunc("/create-user", createUserHandler)
	http.HandleFunc("/get-user/", getUserHandler)
	http.HandleFunc("/update-user/", updateUserHandler)
	http.HandleFunc("/delete-user/", deleteUserHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
