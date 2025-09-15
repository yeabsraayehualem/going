package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	// Register base routes
	app.Router.HandleFunc("/", app.handleHome).Methods("GET")

	// Register app routes
	if err := app.registerAppRoutes(); err != nil {
		log.Printf("Warning: Failed to register app routes: %v", err)
	}
}

// registerAppRoutes finds and registers routes from all apps
func (app *Application) registerAppRoutes() error {
	appsDir := "apps"
	entries, err := os.ReadDir(appsDir)
	if err != nil {
		return fmt.Errorf("error reading apps directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		appName := entry.Name()
		appPath := filepath.Join(appsDir, appName)
		routesFile := filepath.Join(appPath, "routes.go")

		// Check if routes.go exists
		if _, err := os.Stat(routesFile); os.IsNotExist(err) {
			continue
		}

		// Import the app package
		pkgPath := fmt.Sprintf("going/apps/%s", appName)
		appPkg, err := importPackage(pkgPath)
		if err != nil {
			log.Printf("Error importing app %s: %v", appName, err)
			continue
		}

		// Look for RegisterRoutes function
		registerFunc, err := findRegisterRoutesFunc(appPkg, appName)
		if err != nil {
			log.Printf("Error in app %s: %v", appName, err)
			continue
		}

		// Create a subrouter for this app
		router := app.Router.PathPrefix("/" + appName).Subrouter()

		// Call the RegisterRoutes function with the subrouter
		if err := registerFunc(router); err != nil {
			log.Printf("Error registering routes for app %s: %v", appName, err)
			continue
		}

		log.Printf("Registered routes for app: %s", appName)
	}

	return nil
}

// importPackage is a helper to import a package by path
func importPackage(path string) (interface{}, error) {
	// This is a simplified version - in a real implementation, you might use
	// golang.org/x/tools/go/packages or similar to load packages at runtime
	// For now, we'll use a simple approach that works with the existing code
	return nil, nil
}

// findRegisterRoutesFunc looks for a RegisterRoutes function in the package
func findRegisterRoutesFunc(pkg interface{}, appName string) (func(*mux.Router) error, error) {
	// In a real implementation, this would use reflection to find and call the function
	// For now, we'll return a no-op function
	return func(router *mux.Router) error {
		// This will be replaced with actual route registration
		return nil
	}, nil
}

func (app *Application) handleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to going!"))
}
