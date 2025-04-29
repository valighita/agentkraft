package main

import (
	"log"
	"os"
	"valighita/agentkraft/repository"
	configfile "valighita/agentkraft/repository/configFile"
	"valighita/agentkraft/repository/sql"
	"valighita/agentkraft/server"

	"github.com/joho/godotenv"
)

func main() {
	var err error

	if _, err = os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %w", err)
		}
	}

	repoType := os.Getenv("REPO_TYPE")
	if repoType == "" {
		log.Fatal("REPO_TYPE environment variable is not set")
	}

	var agentsRepo repository.AgentsRepository

	switch repoType {
	case "yamlconfig":
		agentsRepo, err = configfile.GetYamlRepositories()
	case "sql":
		agentsRepo, err = sql.GetSqlRepositories()
	default:
		log.Fatalf("Unsupported repository type: %s", repoType)
	}
	if err != nil {
		log.Fatalf("Failed to initialize repositories: %v", err)
	}

	server.NewAgentsHttpServer(agentsRepo).Serve()
}
