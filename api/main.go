package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"skyserver/connection"
	"skyserver/migration"
	"skyserver/models"
	"skyserver/structs"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
)

// type User struct {
// 	ID       int    `json:"id"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

var users = []structs.User{
	{ID: 1, Username: "user1", Password: "password1"},
	{ID: 2, Username: "user2", Password: "password2"},
}

var secretKey = []byte("secret")
var db *sqlx.DB
var migrationFilePath = "./migration.sql"

func main() {
	r := chi.NewRouter()

	db = connection.ConnectToDB()
	defer db.Close()

	if err := migration.ExecuteMigration(db, migrationFilePath); err != nil {
		log.Fatal("Migration failed", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // React development server
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	r.Use(c.Handler)

	r.Post("/api/v1/signup", SignupHandler)
	r.Post("/api/v1/login", LoginHandler)
	r.Post("/api/v1/upload", UploadHandler)
	r.Get("/api/v1/files", ListFilesHandler)

	r.Post("/api/v1/geospatial-data", CreateGeospatialData)
	r.Get("/api/v1/geospatial-data/{id}", GetGeospatialData)
	r.Put("/api/v1/geospatial-data/{id}", UpdateGeospatialData)

	http.ListenAndServe(":8080", r)
}

func CreateGeospatialData(w http.ResponseWriter, r *http.Request) {
	var data structs.GeospatialData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Insert into database
	_, err = db.Exec("INSERT INTO geospatial_data (user_id, geom) VALUES ($1, $2)", data.UserID, data.Geom)
	if err != nil {
		http.Error(w, "Failed to create geospatial data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetGeospatialData(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid geospatial data ID", http.StatusBadRequest)
		return
	}

	var data structs.GeospatialData
	err = db.QueryRow("SELECT id, user_id, geom FROM geospatial_data WHERE id = $1", id).Scan(&data.ID, &data.UserID, &data.Geom)
	if err != nil {
		http.Error(w, "Geospatial data not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func UpdateGeospatialData(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid geospatial data ID", http.StatusBadRequest)
		return
	}

	var data structs.GeospatialData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update database
	_, err = db.Exec("UPDATE geospatial_data SET user_id = $1, geom = $2 WHERE id = $3", data.UserID, data.Geom, id)
	if err != nil {
		http.Error(w, "Failed to update geospatial data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("./uploads")
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	jsonResponse, err := json.Marshal(fileNames)
	if err != nil {
		http.Error(w, "Error marshalling JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonResponse)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // Limit file size to 10 MB

	// Ensure the uploads directory exists
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, "Error creating uploads directory", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := strings.ToLower(handler.Filename)
	ext := filepath.Ext(fileName)

	// Check file extension
	if ext != ".geojson" && ext != ".kml" {
		http.Error(w, "Invalid file format. Only GeoJSON and KML files are allowed", http.StatusBadRequest)
		return
	}

	// Create a new file
	dst, err := os.Create("./uploads/" + fileName)
	if err != nil {
		http.Error(w, "Error creating destination file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully", fileName)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user structs.User

	_ = json.NewDecoder(r.Body).Decode(&user)

	// check if the user already exists
	// userfethced, err := models.GetUserByUsername(db, user.Username)
	// if err != nil {
	// 	log.Println("Username already exists. please choose another username", err)
	// }
	// fmt.Println(userfethced)

	err := models.CreateUser(db, user)
	if err != nil {
		log.Fatal("User not created", err)
	}

	token := generateToken(user.Username)

	json.NewEncoder(w).Encode(map[string]string{"token": token})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	fmt.Println("inside login handler")

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	flag := true

	for _, user := range users {
		if user.Username == credentials.Username && user.Password == credentials.Password {
			token := generateToken(user.Username)
			json.NewEncoder(w).Encode(map[string]string{"token": token})
			flag = false
			break
		}
	}

	if flag {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("User loggedIn"))
}

func generateToken(username string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours
	tokenString, _ := token.SignedString(secretKey)
	return tokenString
}
