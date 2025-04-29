package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"valighita/agentkraft/aiagent"
	"valighita/agentkraft/repository"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type AgentsController struct {
	agentsRepo repository.AgentsRepository
	upgrader   *websocket.Upgrader
}

func NewAgentsController(agentsRepo repository.AgentsRepository) *AgentsController {
	return &AgentsController{
		agentsRepo: agentsRepo,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (c *AgentsController) GetAllAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := c.agentsRepo.GetAllAgents()
	if err != nil {
		http.Error(w, "failad to get all agents", http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]interface{}{
		"agents": agents,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := json.Marshal(jsonResponse)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}

func (c AgentsController) HandleAgentWs(w http.ResponseWriter, r *http.Request) {
	agentIdParam := chi.URLParam(r, "agentID")
	if agentIdParam == "" {
		http.Error(w, "agent ID is required", http.StatusBadRequest)
		return
	}

	agentID, err := strconv.Atoi(agentIdParam)
	if err != nil {
		http.Error(w, "invalid agent ID", http.StatusBadRequest)
		return
	}

	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	agent, err := c.agentsRepo.GetAgentByID(uint(agentID))
	if err != nil || agent == nil {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	aiAgent, err := aiagent.Create(agent)
	if err != nil {
		log.Println("Error creating AI agent:", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		response, err := aiAgent.GetCompletion(string(msg))
		if err != nil {
			response = fmt.Sprintf("Error getting response: %v", err)
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(strings.TrimSpace(response)))
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
