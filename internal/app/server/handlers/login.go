package handlers

import (
	"Forum/internal/model"
	"Forum/internal/store"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

type RegisterResponse struct {
	Error string `json:"error,omitempty"`
}

type LoginResponse struct {
	Error string `json:"error,omitempty"`
}

func Register(store store.Store, templates *template.Template, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var data struct {
				UserName   string `json:"userName"`
				Email      string `json:"email"`
				Password   string `json:"password"`
				RePassword string `json:"rePassword"`
				FirstName  string `json:"firstName"`
				LastName   string `json:"lastName"`
				Age        int    `json:"age"`
				Gender     string `json:"gender"`
			}

			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(RegisterResponse{Error: "Invalid request payload"})
				return
			}

			if data.Password != data.RePassword {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(RegisterResponse{Error: "Passwords don't match"})
				return
			}

			err := store.User().ExistingUser(data.UserName, data.Email)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(RegisterResponse{Error: "User already exists"})
				return
			}

			user, err := model.NewUser(data.UserName, data.Email, data.Password, data.FirstName, data.LastName, data.Age, data.Gender)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(RegisterResponse{Error: "Failed to create the user"})
				return
			}

			if err = store.User().Register(user); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(RegisterResponse{Error: "Failed to register the user"})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(RegisterResponse{})
		}
	}
}

func Login(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var data struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}

			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(LoginResponse{Error: "Invalid request payload"})
				return
			}

			user := &model.User{
				Email:    data.Email,
				Password: data.Password,
			}

			if err := store.User().Login(user); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(LoginResponse{Error: "Invalid email or password"})
				return
			}

			session, err := model.NewSession(user.UUID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(LoginResponse{Error: "Failed to create session"})
				return
			}

			if err := store.Session().Create(session); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(LoginResponse{Error: "Failed to create session"})
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_uuid",
				Value:    session.SessionID,
				Expires:  session.ExpiresAt,
				HttpOnly: true,
				Secure:   false, // Set to true if you have HTTPS
			})

			logger.Printf("Session created: %v\n", session)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(LoginResponse{})
		}
	}
}

func Logout(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_uuid")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := store.Session().Delete(cookie.Value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_uuid",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   false, // Set to true if you have HTTPS
		})

		logger.Printf("Session deleted: %s\n", cookie.Value)

		w.WriteHeader(http.StatusOK)
	}
}

func GetUsersHandler(store store.Store, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := store.User().GetUsersOrderedByLastMessageOrAlphabetically()
		if err != nil {
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
