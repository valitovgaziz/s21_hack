package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"api_hr/internal/handlers"
	"api_hr/internal/middleware"
	"api_hr/internal/repository"
	"api_hr/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Server struct {
	router *chi.Mux
	db     *gorm.DB
}

func New(db *gorm.DB) *Server {
	s := &Server{
		router: chi.NewRouter(),
		db:     db,
	}
	s.configureRouter(db)
	return s
}

func (s *Server) configureRouter(db *gorm.DB) {
	for _, middleware := range handlers.CommonMiddleware() {
		s.router.Use(middleware)
	}

	s.router.Get("/health", s.healthCheck)

	s.router.Route("/v1", func(r chi.Router) {
		r.Get("/check", s.healthCheck)
		s.setupUserRoutes(r, db)
		s.setupSurveyRoutes(r, db)
		s.setupResponseRoutes(r, db)
		s.setupNotificationRoutes(r, db)
		s.setupAnalyticsRoutes(r, db)
	})

	chi.Walk(s.router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s] %s\n", method, route)
		return nil
	})
}

func (s *Server) setupUserRoutes(r chi.Router, db *gorm.DB) {
	userRepo := repository.NewUserRepository(s.db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := &handlers.AuthHandler{DB: db}

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/login-phone", authHandler.LoginPhone)
		r.Post("/verify-otp", authHandler.VerifyOTP)
		r.Get("/validate", authHandler.ValidateToken)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.GetAllUsers)
			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUser)
		})
	})
}

func (s *Server) setupSurveyRoutes(r chi.Router, db *gorm.DB) {
	surveyRepo := repository.NewSurveyRepository(s.db)
	responseRepo := repository.NewResponseRepository(s.db)
	surveyService := service.NewSurveyService(surveyRepo, responseRepo)
	surveyHandler := handlers.NewSurveyHandler(surveyService)

	r.Route("/surveys", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Post("/", surveyHandler.Create)
			r.Put("/{id}", surveyHandler.Update)
			r.Delete("/{id}", surveyHandler.Delete)
			r.Post("/{id}/publish", surveyHandler.Publish)
			r.Post("/{id}/complete", surveyHandler.Complete)
			r.Post("/{id}/archive", surveyHandler.Archive)
			r.Get("/{id}/stats", surveyHandler.GetStats)
		})

		r.Get("/", surveyHandler.GetAll)
		r.Get("/{id}", surveyHandler.GetByID)
	})
}

func (s *Server) setupResponseRoutes(r chi.Router, db *gorm.DB) {
	responseRepo := repository.NewResponseRepository(s.db)
	surveyRepo := repository.NewSurveyRepository(s.db)
	responseService := service.NewResponseService(responseRepo, surveyRepo)
	responseHandler := handlers.NewResponseHandler(responseService)

	r.Route("/responses", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/{surveyId}/start", responseHandler.StartSurvey)
		r.Post("/{responseId}/answer", responseHandler.SubmitAnswer)
		r.Post("/{responseId}/complete", responseHandler.CompleteSurvey)
		r.Get("/{id}", responseHandler.GetResponse)
		r.Get("/survey/{surveyId}", responseHandler.GetSurveyResponses)
		r.Get("/survey/{surveyId}/check", responseHandler.CheckResponded)
	})
}

func (s *Server) setupNotificationRoutes(r chi.Router, db *gorm.DB) {
	notifRepo := repository.NewNotificationRepository(s.db)
	notifService := service.NewNotificationService(notifRepo)
	notifHandler := handlers.NewNotificationHandler(notifService)

	r.Route("/notifications", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/subscribe", notifHandler.Subscribe)
		r.Get("/subscriptions", notifHandler.GetSubscriptions)
		r.Delete("/subscriptions/{id}", notifHandler.Unsubscribe)
		r.Get("/preferences", notifHandler.GetPreferences)
		r.Put("/preferences", notifHandler.UpdatePreferences)
		r.Get("/logs/{surveyId}", notifHandler.GetSurveyLogs)
	})
}

func (s *Server) setupAnalyticsRoutes(r chi.Router, db *gorm.DB) {
	surveyRepo := repository.NewSurveyRepository(s.db)
	responseRepo := repository.NewResponseRepository(s.db)
	notifRepo := repository.NewNotificationRepository(s.db)
	analyticsService := service.NewAnalyticsService(surveyRepo, responseRepo, notifRepo)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	r.Route("/analytics", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/surveys/{surveyId}", analyticsHandler.GetSurveyAnalytics)
	})
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := s.db.DB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusServiceUnavailable)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		http.Error(w, "Database ping failed", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC1123),
	})
}

func (s *Server) Run(port string) error {
	return http.ListenAndServe(":"+port, s.router)
}
