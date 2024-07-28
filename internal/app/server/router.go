package server

import (
	"Forum/internal/app/server/handlers"
)

func (s *Server) HandlePaths() {
	s.router.Handle("/static/", handlers.ServeStatic())
	s.router.HandleFunc("/", handlers.Home(s.Store, s.Templates, s.Logger))
	s.router.HandleFunc("/registerPage", handlers.RegisterPage(s.Templates))
	s.router.HandleFunc("/saveUser", handlers.SaveRegister(s.Store, s.Templates, s.Logger))
	s.router.HandleFunc("/loginPage", handlers.LoginPage(s.Templates))
	s.router.HandleFunc("/login", handlers.Login(s.Store, s.Templates, s.Logger))
	s.router.HandleFunc("/createPost", handlers.CreatePost(s.Store, s.Logger))
	s.router.HandleFunc("/createPostPage", handlers.CreatePostPage(s.Store, s.Templates))
	s.router.HandleFunc("/createCategory", handlers.CreateCategory(s.Store, s.Templates))
	s.router.HandleFunc("/createCategoryPage", handlers.CreateCategoryPage(s.Store, s.Templates))
	s.router.HandleFunc("/category/", handlers.CategoryPosts(s.Store, s.Templates))
	s.router.HandleFunc("/userProfilePage", handlers.ServeUserProfile(s.Store, s.Templates))
	s.router.HandleFunc("/logout", handlers.Logout(s.Store, s.Logger))
	s.router.HandleFunc("/createComment", handlers.CreateComment(s.Store, s.Templates, s.Logger))
	s.router.HandleFunc("/createPostReaction", handlers.HandleCreatePostReaction(s.Store, s.Logger))
	s.router.HandleFunc("/reactComment", handlers.HandleCreateCommentReaction(s.Store, s.Logger))
	s.router.HandleFunc("/send-message", handlers.SendMessageHandler(s.Store, s.Logger))
	s.router.HandleFunc("/get-messages", handlers.GetMessagesHandler(s.Store, s.Logger))
	s.router.HandleFunc("/update-status", handlers.UpdateStatusHandler(s.Store, s.Logger))
	s.router.HandleFunc("/get-user-status", handlers.GetUserStatusHandler(s.Store, s.Logger))
	s.router.HandleFunc("/mark-message-read", handlers.MarkMessageAsReadHandler(s.Store, s.Logger))
	s.router.HandleFunc("/get-users", handlers.GetUsersHandler(s.Store, s.Logger))
	s.router.HandleFunc("/ws", handlers.HandleConnections(s.Store))
	go handlers.HandleMessages(s.Store)
}
