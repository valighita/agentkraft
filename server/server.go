package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"valighita/agentkraft/repository"
	"valighita/agentkraft/server/controller"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	defaultServerPort = 8080
)

type AgentsHttpServer struct {
	agentsRepo repository.AgentsRepository
	router     *chi.Mux
	listenAddr string
	port       int
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewAgentsHttpServer(agentsRepo repository.AgentsRepository) *AgentsHttpServer {
	serverPort := defaultServerPort
	listenAddr := "127.0.0.1"

	if addr, exists := os.LookupEnv("HTTP_SERVER_LISTEN_ADDR"); exists {
		listenAddr = addr
	}

	if envVar, exists := os.LookupEnv("HTTP_SERVER_PORT"); exists {
		var err error
		serverPort, err = strconv.Atoi(envVar)
		if err != nil {
			log.Fatalf("Error parsing HTTP_SERVER_PORT: %s", err)
		}
	}

	server := AgentsHttpServer{
		agentsRepo: agentsRepo,
		listenAddr: listenAddr,
		port:       serverPort,
	}

	server.setup()

	return &server
}

func (server *AgentsHttpServer) setup() {
	r := chi.NewRouter()

	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	agentsController := controller.NewAgentsController(server.agentsRepo)

	r.Route("/agents", func(r chi.Router) {
		r.Get("/", agentsController.GetAllAgents)
		r.Get("/ws/{agentID}", agentsController.HandleAgentWs)
	})

	// serve frontend/index.html on "/" path
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "server/frontend/index.html")
	})

	server.router = r
}

func (server AgentsHttpServer) Serve() {
	log.Printf("Starting server on %s:%d", server.listenAddr, server.port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", server.listenAddr, server.port), server.router)
	if err != nil {
		log.Fatal(err)
	}
}
