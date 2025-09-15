package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	app "going/internal/app"
	"going/internal/config"
	"github.com/gorilla/mux"
)

const (
	configPath = "config/config.yaml"
)

func main() {
	// Command line flags
	initFlag := flag.Bool("init", false, "Initialize a new going project")
	createAppFlag := flag.String("create-app", "", "Create a new app with the given name")
	flag.Parse()

	switch {
	case *initFlag:
		if err := initializeProject(); err != nil {
			log.Fatalf("Failed to initialize project: %v", err)
		}
		fmt.Println("Project initialized successfully!")
		return
	case *createAppFlag != "":
		if err := createApp(*createAppFlag); err != nil {
			log.Fatalf("Failed to create app: %v", err)
		}
		fmt.Printf("App '%s' created successfully!\n", *createAppFlag)
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize and start the application
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

// initializeProject sets up a new going project structure
// createApp creates a new app with the given name
func createApp(appName string) error {
	// Create app directory
	appDir := filepath.Join("apps", appName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("error creating app directory: %w", err)
	}

	// Create models.go file
	modelsContent := fmt.Sprintf(`package %s

import (
	"going/internal/database"
)

type %sModel struct {
	ID   uint   ` + "`" + `gorm:"primaryKey"` + "`" + `
	Name string ` + "`" + `gorm:"size:255"` + "`" + `
}

func init() {
	// Register your models here
	database.RegisterModels(&%sModel{})
}
`, appName, strings.Title(appName), strings.Title(appName))

	modelsPath := filepath.Join(appDir, "models.go")
	if err := os.WriteFile(modelsPath, []byte(modelsContent), 0644); err != nil {
		return fmt.Errorf("error creating models file: %w", err)
	}

	// Create routes.go file
	routesContent := fmt.Sprintf(`package %s

import (
	"net/http"
)

// RegisterRoutes registers all routes for this app
func RegisterRoutes(router *mux.Router) {
	// Register your routes here
	// Example:
	// router.HandleFunc("/%s", handle%s).Methods("GET")
}
`, appName, appName, strings.Title(appName))

	routesPath := filepath.Join(appDir, "routes.go")
	if err := os.WriteFile(routesPath, []byte(routesContent), 0644); err != nil {
		return fmt.Errorf("error creating routes file: %w", err)
	}

	return nil
}

func initializeProject() error {
	// Create necessary directories
	dirs := []string{
		"apps",
		"migrations",
		"templates",
		"static/css",
		"static/js",
		"static/images",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	// Create a default config if it doesn't exist
	cfgPath := filepath.Join("config", "config.yaml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		cfg := config.DefaultConfig()
		if err := cfg.Save(cfgPath); err != nil {
			return fmt.Errorf("error creating default config: %w", err)
		}
	}

	// Create a sample app
	if err := os.MkdirAll("apps/example", 0755); err != nil {
		return fmt.Errorf("error creating example app: %w", err)
	}

	// Create a sample app file
	sampleAppContent := `package example

import (
	"going/internal/database"
)

type ExampleModel struct {
	ID   uint   ` + "`" + `gorm:"primaryKey"` + "`" + `
	Name string ` + "`" + `gorm:"size:255"` + "`" + `
}

func init() {
	// Register your models here
	database.RegisterModels(&ExampleModel{})
}
`

	if err := os.WriteFile("apps/example/models.go", []byte(sampleAppContent), 0644); err != nil {
		return fmt.Errorf("error creating example app file: %w", err)
	}

	return nil
}
