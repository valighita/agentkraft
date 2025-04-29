package aiagent

import (
	"context"
	"errors"
	"log"
	"time"
	"valighita/agentkraft/repository"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
)

const (
	defaultMaxTurns       = 10
	contextPromptTemplate = "It's important to only answer relevant questions about the services provided, do not provide information about unrelated topics.\n\n" +
		"{{.tool_descriptions}}"
)

type aiAgentFactory interface {
	createLLM(llmModel string, apiKey string) (llms.Model, error)
}

type openAIAgentFactory struct {
}

type AiAgent struct {
	executor    *agents.Executor
	temperature float64
}

func (f *openAIAgentFactory) createLLM(llmModel string, apiKey string) (llms.Model, error) {
	return openai.New(
		openai.WithModel(llmModel),
		openai.WithToken(apiKey),
	)
}

func getAgentFactory(modelProvider repository.AIModelProvider) aiAgentFactory {
	switch modelProvider {
	case repository.AIModelProviderOpenAI:
		return &openAIAgentFactory{}
	default:
		return nil
	}
}

func Create(agent *repository.Agent) (*AiAgent, error) {
	memory := memory.NewConversationBuffer()

	tools := GetAgentTools(agent)
	aiFactory := getAgentFactory(agent.AIModelProvider)
	if aiFactory == nil {
		return nil, errors.New("unsupported AI model provider")
	}

	llm, err := aiFactory.createLLM(agent.AIModelType, agent.AIModelToken)
	if err != nil {
		log.Fatalf("Error creating LLM: %v", err)
	}

	aiAgent := agents.NewConversationalAgent(llm,
		tools,
		agents.WithPromptPrefix(agent.SystemPrompt+"\n\n"+
			contextPromptTemplate+"\n\nCurrent time is "+time.Now().Format("2006-01-02 15:04:05, Monday")),
		agents.WithMemory(memory),
	)

	executor := agents.NewExecutor(aiAgent,
		agents.WithMaxIterations(agent.MaxTurns),
		agents.WithMemory(memory),
	)

	return &AiAgent{
		executor:    executor,
		temperature: agent.AIModelTemperature,
	}, nil
}

func (a *AiAgent) GetCompletion(prompt string) (string, error) {
	return chains.Run(context.Background(), a.executor, prompt,
		chains.WithTemperature(a.temperature),
	)
}
