package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Estructura para los datos del usuario
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Respuesta estándar de la API
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Base de datos en memoria (en producción usarías una DB real)
var users = []User{
	{ID: 1, Name: "Juan Pérez", Email: "juan@example.com"},
	{ID: 2, Name: "María García", Email: "maria@example.com"},
}

var nextID = 3

// Middleware para logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		
		next.ServeHTTP(w, r)
		
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// Middleware para CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// Health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Status:  "success",
		Message: "Service is healthy",
		Data: map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    time.Since(startTime).String(),
		},
	}
	json.NewEncoder(w).Encode(response)
}

var startTime = time.Now()

// Obtener todos los usuarios
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data:    users,
	}
	json.NewEncoder(w).Encode(response)
}

// Obtener un usuario por ID
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Status:  "error",
			Message: "Invalid user ID",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	for _, user := range users {
		if user.ID == id {
			response := Response{
				Status:  "success",
				Message: "User found",
				Data:    user,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	
	w.WriteHeader(http.StatusNotFound)
	response := Response{
		Status:  "error",
		Message: "User not found",
	}
	json.NewEncoder(w).Encode(response)
}

// Crear un nuevo usuario
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Status:  "error",
			Message: "Invalid JSON format",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Validación básica
	if newUser.Name == "" || newUser.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Status:  "error",
			Message: "Name and email are required",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Asignar ID y agregar a la lista
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)
	
	w.WriteHeader(http.StatusCreated)
	response := Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    newUser,
	}
	json.NewEncoder(w).Encode(response)
}

// Actualizar un usuario
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Status:  "error",
			Message: "Invalid user ID",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Status:  "error",
			Message: "Invalid JSON format",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	for i, user := range users {
		if user.ID == id {
			updatedUser.ID = id
			users[i] = updatedUser
			response := Response{
				Status:  "success",
				Message: "User updated successfully",
				Data:    updatedUser,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	
	w.WriteHeader(http.StatusNotFound)
	response := Response{
		Status:  "error",
		Message: "User not found",
	}
	json.NewEncoder(w).Encode(response)
}

// Eliminar un usuario
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Status:  "error",
			Message: "Invalid user ID",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			response := Response{
				Status:  "success",
				Message: "User deleted successfully",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	
	w.WriteHeader(http.StatusNotFound)
	response := Response{
		Status:  "error",
		Message: "User not found",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Crear router
	r := mux.NewRouter()
	
	// Aplicar middlewares
	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)
	
	// Definir rutas
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/users", getUsersHandler).Methods("GET")
	r.HandleFunc("/api/users/{id}", getUserHandler).Methods("GET")
	r.HandleFunc("/api/users", createUserHandler).Methods("POST")
	r.HandleFunc("/api/users/{id}", updateUserHandler).Methods("PUT")
	r.HandleFunc("/api/users/{id}", deleteUserHandler).Methods("DELETE")
	
	// Configurar puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)
	log.Printf("API endpoints available at: http://localhost:%s/api/users", port)
	
	// Iniciar servidor
	log.Fatal(http.ListenAndServe(":"+port, r))
}