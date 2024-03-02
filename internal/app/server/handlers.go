package server

import (
	"Forum/internal/model"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var templates = template.Must(template.ParseGlob("./web/templates/*.html"))

func (s *server) HandlePaths() {
	s.router.Handle("/static/", s.serveStatic())
	s.router.HandleFunc("/", s.home())
	s.router.HandleFunc("/registerPage", s.registerPage())
	s.router.HandleFunc("/saveUser", s.saveRegister())
	s.router.HandleFunc("/loginPage", s.loginPage())
	s.router.HandleFunc("/login", s.login())
	s.router.HandleFunc("/createPost", s.createPost())
	s.router.HandleFunc("/createPostPage", s.createPostPage())
	s.router.HandleFunc("/createCategory", s.createCategory())
	s.router.HandleFunc("/createCategoryPage", s.createCategoryPage())
	s.router.HandleFunc("/category/", s.categoryPosts())
}

func (s *server) registerPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		execTmpl(w, templates.Lookup("registerPage.html"), nil)
	}
}

func (s *server) saveRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userName := r.FormValue("userName")
		email := r.FormValue("email")
		password := r.FormValue("password")

		err := s.store.User().ExistingUser(userName, email)
		if err != nil {
			http.Redirect(w, r, "/registerPage", http.StatusSeeOther)
			s.logger.Println("redirect - error:", err)
			return
		}

		user, err := model.NewUser(userName, email, password)
		if err != nil {
			s.logger.Println("NewUser() error: ", err)
			http.Redirect(w, r, "/registerPage", http.StatusSeeOther)
			return
		}

		if err = s.store.User().Register(user); err != nil {
			s.logger.Println("Register() error: ", err)
			http.Redirect(w, r, "/registerPage", http.StatusSeeOther)
			return
		}

		execTmpl(w, templates.Lookup("main.html"), nil)

	}
}

func (s *server) loginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		execTmpl(w, templates.Lookup("login.html"), nil)
	}
}

func (s *server) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Println("@ login page")
		email := r.FormValue("email")
		password := r.FormValue("password")
		s.logger.Println(email, password)

		// Create a new User instance with only email and password
		user := &model.User{
			Email:    email,
			Password: password,
		}

		if err := s.store.User().Login(user); err != nil {
			s.logger.Println("redirect - Login() error: ", err)
			http.Redirect(w, r, "/loginPage", http.StatusBadRequest)
			return
		}

		//	execTmpl(w, templates.Lookup("login.html"), nil)
	}

}

func (s *server) serveStatic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))).ServeHTTP(w, r)
	}
}

func (s *server) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetching all posts.
		posts, err := s.store.Post().GetAll()
		if err != nil {
			s.logger.Println("error fetching posts:", err)
			http.Error(w, "error fetching posts", http.StatusInternalServerError)
			return
		}

		// Fetching categories for each post.
		for _, post := range posts {
			categories, err := s.store.Post().GetCategories(post.ID)
			if err != nil {
				s.logger.Println("error fetching categories for post:", err)
				http.Error(w, "error fetching post categories", http.StatusInternalServerError)
				return
			}
			post.Categories = categories
		}

		// Fetching all categories.
		allCategories, err := s.store.Category().GetAll()
		if err != nil {
			s.logger.Println("error fetching categories:", err)
			http.Error(w, "error fetching categories", http.StatusInternalServerError)
			return
		}

		// Struct to pass into template.
		data := struct {
			Posts      []*model.Post
			Categories []*model.Category
		}{
			Posts:      posts,
			Categories: allCategories,
		}

		execTmpl(w, templates.Lookup("main.html"), data)
	}
}

// execTmpl renders a template with the given data or returns an internal server error.
func execTmpl(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)
	}
}

func (s *server) createPostPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := s.store.Category().GetAll()
		if err != nil {
			// handle error
		}

		data := struct {
			Categories []*model.Category
		}{
			Categories: categories,
		}

		execTmpl(w, templates.Lookup("createPostPage.html"), data)
	}
}

func (s *server) createPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := "Admin"
		subject := r.FormValue("postTitle")
		content := r.FormValue("postText")

		post, err := model.NewPost(userID, subject, content)
		if err != nil {
			s.logger.Println("NewPost() error: ", err)
			http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
			return
		}

		if err = s.store.Post().Create(post); err != nil {
			s.logger.Println("Create() error: ", err)
			http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
			return
		}

		// Parse form for category checkboxes
		r.ParseForm()
		categoryIDs := r.PostForm["categoryIDs"]

		for _, categoryIDStr := range categoryIDs {
			categoryID, err := strconv.Atoi(categoryIDStr)
			if err != nil {
				s.logger.Println("Error converting categoryID to int: ", err)
				http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
				return
			}

			if err := s.store.Post().AddCategoryToPost(post.ID, categoryID); err != nil {
				s.logger.Println("Error adding category to post: ", err)
				http.Redirect(w, r, "/createPostPage", http.StatusSeeOther)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *server) createCategoryPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		execTmpl(w, templates.Lookup("createCategory.html"), nil)
	}
}

func (s *server) createCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryName := r.FormValue("categoryName")
		category := &model.Category{Name: categoryName}

		if err := s.store.Category().Create(category); err != nil {
			s.logger.Println("Create category error: ", err)
			http.Redirect(w, r, "/createCategoryPage", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *server) categoryPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Pull out the categoryID from the url.
		s.logger.Println("Path:", r.URL.Path)

		categoryIDStr := strings.TrimPrefix(r.URL.Path, "/category/")
		s.logger.Println("Parsed CategoryIDStr:", categoryIDStr)

		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			s.logger.Println("Error converting category ID:", err)
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		s.logger.Println("Using CategoryID:", categoryID)

		// Get all posts in the category.
		posts, err := s.store.Post().GetByCategory(categoryID)
		if err != nil {
			s.logger.Println("Error fetching posts:", err)
			http.Error(w, "error fetching posts by category", http.StatusInternalServerError)
			return
		}

		for _, post := range posts {
			s.logger.Println("Post subject:", post.Subject)
		}

		allCategories, err := s.store.Category().GetAll()
		if err != nil {
			s.logger.Println("Error fetching categories:", err)
			http.Error(w, "error fetching all categories", http.StatusInternalServerError)
			return
		}

		for _, category := range allCategories {
			s.logger.Println("Category name:", category.Name)
		}

		data := struct {
			Posts      []*model.Post
			Categories []*model.Category
		}{
			Posts:      posts,
			Categories: allCategories,
		}

		execTmpl(w, templates.Lookup("home.html"), data)
	}
}
