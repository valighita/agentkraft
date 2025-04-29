package repository

import (
	"time"
)

type AIModelProvider string

const (
	AIModelProviderOpenAI AIModelProvider = "openai"
)

var AIModelProviders = []AIModelProvider{
	AIModelProviderOpenAI,
}

type HttpToolHeader struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// HttpTool represents a HTTP tool that can be used by a agent.
type HttpTool struct {
	Name        string           `json:"name" yaml:"name"`
	Description string           `json:"description" yaml:"description"`
	Url         string           `json:"url" yaml:"url"`
	HttpMethod  string           `json:"method" yaml:"method"`
	Headers     []HttpToolHeader `json:"headers" yaml:"headers"`
	Params      []string         `json:"params" yaml:"params"`
}

// Agent holds the configuration for a agent.
type Agent struct {
	ID                  uint            `gorm:"primaryKey" yaml:"id" json:"id"`
	Name                string          `gorm:"size:255;not null" yaml:"name" json:"name"`
	SystemPrompt        string          `gorm:"size:255;not null" yaml:"systemPrompt" json:"systemPrompt"`   // System prompt for the agent
	AIModelProvider     AIModelProvider `gorm:"size:255;not null" yaml:"llmProvider" json:"llmProvider"`     // Provider of the AI model
	AIModelType         string          `gorm:"size:255;not null" yaml:"llmModel" json:"llmModel"`           // Type of the AI model
	AIModelToken        string          `gorm:"size:255;not null" yaml:"-" json:"-"`                         // API token for accessing the AI model
	AIModelApiKeyEnvVar string          `gorm:"size:255;not null" yaml:"llmApiKeyEnvVar" json:"-"`           // Environment variable for the API key
	AIModelTemperature  float64         `gorm:"default:0.7" yaml:"llmTemperature" json:"llmTemperature"`     // Temperature for the AI model
	MaxTurns            int             `gorm:"default:10" yaml:"maxTurns" json:"maxTurns"`                  // Maximum number of turns for the agent
	HttpTools           []HttpTool      `gorm:"type:text;serializer:json" yaml:"httpTools" json:"httpTools"` // HTTP tools for the agent, stored as JSON encoded string in the SQL db
	CreatedAt           time.Time       `gorm:"autoCreateTime"`
	UpdatedAt           time.Time       `gorm:"autoUpdateTime"`
}

type AgentsRepository interface {
	GetAgentByID(id uint) (*Agent, error)
	GetAllAgents() ([]Agent, error)
}
