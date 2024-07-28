package handlers

import (
	"Forum/internal/model"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"Forum/internal/store"
)

// RegisterPage обробляє відображення сторінки реєстрації
func RegisterPage(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		execTmpl(w, templates.Lookup("registerPage.html"), nil)
	}
}

// SaveRegister обробляє збереження реєстраційних даних
func SaveRegister(store store.Store, templates *template.Template, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := model.RegisterPageData{}

		userName := r.FormValue("userName")
		email := r.FormValue("email")
		password := r.FormValue("password")
		rePassword := r.FormValue("rePassword")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		age, _ := strconv.Atoi(r.FormValue("age"))
		gender := r.FormValue("gender")

		// Перевірка чи співпадають паролі
		if password != rePassword {
			logger.Println("Passwords don't match")
			data.ErrorMsg = "Passwords don't match"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		err := store.User().ExistingUser(userName, email)
		if err != nil {
			logger.Println("error:", err)
			data.UserExistsErrorMsg = "User already exists in the system"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		user, err := model.NewUser(userName, email, password, firstName, lastName, age, gender)
		if err != nil {
			logger.Println("NewUser() error: ", err)
			data.ErrorMsg = "Failed to create the user"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		if err = store.User().Register(user); err != nil {
			logger.Println("Register() error: ", err)
			data.ErrorMsg = "Failed to register the user"
			execTmpl(w, templates.Lookup("registerPage.html"), data)
			return
		}

		execTmpl(w, templates.Lookup("main.html"), nil)
	}
}

// LoginPage обробляє відображення сторінки входу
func LoginPage(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorMessage := ""
		if errMsg := r.URL.Query().Get("error"); errMsg == "notfound" {
			errorMessage = "User not found. Please try again."
		}

		execTmpl(w, templates.Lookup("login.html"), map[string]string{"ErrorMessage": errorMessage})
	}
}

// Login обробляє вхід користувача
func Login(store store.Store, templates *template.Template, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println("@ login page")
		email := r.FormValue("email")
		password := r.FormValue("password")
		logger.Println(email)

		user := &model.User{
			Email:    email,
			Password: password,
		}

		// Аутентифікація користувача
		if err := store.User().Login(user); err != nil {
			logger.Println("Login() error: ", err)
			http.Redirect(w, r, "/loginPage?error=notfound", http.StatusSeeOther)
			return
		}

		// Створення нової сесії для користувача
		session, err := model.NewSession(user.UUID)
		if err != nil {
			logger.Println("NewSession() error: ", err)
			http.Redirect(w, r, "/loginPage", http.StatusInternalServerError)
			return
		}

		// Збереження сесії в базі даних
		if err := store.Session().Create(session); err != nil {
			logger.Println("CreateSession() error: ", err)
			http.Redirect(w, r, "/loginPage", http.StatusInternalServerError)
			return
		}

		// Встановлення cookie сесії
		http.SetCookie(w, &http.Cookie{
			Name:     "session_uuid",
			Value:    session.SessionID,
			Expires:  session.ExpiresAt,
			HttpOnly: true,
			Secure:   false, // Встановити true, якщо використовуєте HTTPS
		})

		// Перенаправлення користувача на головну сторінку
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

// ServeUserProfile обробляє відображення профілю користувача
func ServeUserProfile(store store.Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Отримання cookie сесії
		sessionCookie, err := r.Cookie("session_uuid")
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Отримання сесії
		session, err := store.Session().GetByUUID(sessionCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Отримання користувача
		user, err := store.User().GetByUUID(session.UserUUID)
		if err != nil {
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			return
		}

		// Рендеринг шаблону з даними користувача
		execTmpl(w, templates.Lookup("userProfilePage.html"), user)
	}
}

// GetUsersHandler обробляє отримання списку користувачів
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
