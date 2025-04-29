package configfile

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"valighita/agentkraft/repository"

	"gopkg.in/yaml.v2"
)

type YAMLAgentsRepository struct {
	filePath string
	agents   []repository.Agent
}

type yamlData struct {
	Agents []repository.Agent `yaml:"agents"`
}

func NewYAMLAgentsRepository() (*YAMLAgentsRepository, error) {
	filePath := os.Getenv("YAML_CONFIG_FILE")
	if filePath == "" {
		return nil, errors.New("YAML_CONFIG_FILE environment variable is not set")
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	agents, err := loadAgentsFromConfigFile(absPath)
	if err != nil {
		return nil, err
	}

	return &YAMLAgentsRepository{
		filePath: absPath,
		agents:   agents,
	}, nil
}

func validateHttpTool(tool *repository.HttpTool) error {
	if tool == nil {
		return errors.New("tool is nil")
	}
	if tool.Name == "" {
		return errors.New("tool name is empty")
	}
	if tool.Url == "" {
		return errors.New("tool URL is empty")
	}
	if _, err := url.Parse(tool.Url); err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	if tool.HttpMethod == "" {
		tool.HttpMethod = "GET"
	}
	for _, param := range tool.Params {
		if param == "" {
			return errors.New("tool parameter is empty")
		}
	}
	for _, header := range tool.Headers {
		if header.Key == "" {
			return errors.New("tool header key is empty")
		}
	}
	return nil
}

func loadAgentsFromConfigFile(filePath string) ([]repository.Agent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data yamlData
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	for i, agent := range data.Agents {
		if agent.ID == 0 {
			agent.ID = uint(i + 1)
		}

		if agent.Name == "" {
			agent.Name = "Agent " + string(i+1)
		}

		if agent.AIModelProvider == "" || agent.AIModelType == "" || agent.AIModelApiKeyEnvVar == "" {
			return nil, fmt.Errorf("agent %d: attributes llmProvider, llmType and llmApiKey myst be specified", agent.ID)
		}

		apiKey := os.Getenv(agent.AIModelApiKeyEnvVar)
		if apiKey == "" {
			return nil, fmt.Errorf("agent %d: environment variable %s is not set", agent.ID, agent.AIModelApiKeyEnvVar)
		}
		agent.AIModelToken = apiKey

		for j, tool := range agent.HttpTools {
			if err := validateHttpTool(&tool); err != nil {
				return nil, fmt.Errorf("agent %d, tool %d: %v", agent.ID, j+1, err)
			}
		}

		data.Agents[i] = agent
	}

	return data.Agents, nil
}

func (r *YAMLAgentsRepository) GetAgentByID(id uint) (*repository.Agent, error) {
	for _, agent := range r.agents {
		if agent.ID == id {
			return &agent, nil
		}
	}
	return nil, nil
}

func (r *YAMLAgentsRepository) GetAllAgents() ([]repository.Agent, error) {
	return r.agents, nil
}
