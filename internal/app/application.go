package app

import (
	"database/sql"
	"log"
	"net/http"

	"going/internal/config"
	"going/internal/database"
	"going/internal/middleware"
	"going/internal/session"

	"github.com/gorilla/mux"
)

type Application struct {
	Config  *config.Config
	DB      *sql.DB
	Router  *mux.Router
	Session *session.Manager
}

func NewApplication(cfg *config.Config) (*Application, error) {
	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize session manager
	sessionManager := session.NewManager(cfg)

	// Create router
	router := mux.NewRouter()

	return &Application{
		Config:  cfg,
		DB:      db,
		Router:  router,
		Session: sessionManager,
	}, nil
}

func (app *Application) Run() error {
	// Register routes
	app.registerRoutes()

	// Create a new router with the logging middleware
	loggedRouter := middleware.LoggingMiddleware(app.Router)

	// Start the server
	serverAddr := app.Config.Server.Host + ":" + app.Config.Server.Port
	log.Printf("Starting server on %s\n", serverAddr)
	return http.ListenAndServe(serverAddr, loggedRouter)
}

func (app *Application) registerRoutes() {
	// Register your routes here
	app.Router.HandleFunc("/", app.handleHome).Methods("GET")
}

func (app *Application) handleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to going!"))
}
