package aiagent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"valighita/agentkraft/repository"

	langchaintools "github.com/tmc/langchaingo/tools"
)

type agentTool struct {
	info repository.HttpTool
}

func (t *agentTool) Name() string {
	return t.info.Name
}

func (t *agentTool) Description() string {
	return t.info.Description + "\n" +
		"Params: " + strings.Join(t.info.Params, ", ") + "\n" +
		"Parameters must be passed as a JSON object"
}

func (t *agentTool) Call(ctx context.Context, input string) (string, error) {
	var inputMap map[string]string
	err := json.Unmarshal([]byte(input), &inputMap)
	if err != nil {
		return fmt.Sprintf("failed to unmarshal input: %v", err), nil
	}

	reqUrl := t.info.Url

	paramsMap := map[string]string{}
	for _, arg := range t.info.Params {
		if _, ok := inputMap[arg]; !ok || inputMap[arg] == "" {
			return fmt.Sprintf("missing argument: %s", arg), nil
		}

		if strings.Contains(reqUrl, fmt.Sprintf("{%s}", arg)) {
			reqUrl = strings.Replace(reqUrl, fmt.Sprintf("{%s}", arg), url.QueryEscape(inputMap[arg]), -1)
			continue
		}

		paramsMap[arg] = inputMap[arg]
	}

	var body io.Reader
	if t.info.HttpMethod == "GET" {
		newUrl := reqUrl
		if !strings.Contains(reqUrl, "?") {
			newUrl += "?"
		} else {
			newUrl += "&"
		}
		for k, v := range paramsMap {
			newUrl += fmt.Sprintf("%s=%s&", k, url.QueryEscape(v))
		}
		reqUrl = strings.TrimSuffix(newUrl, "&")
	} else {
		data, err := json.Marshal(paramsMap)
		if err != nil {
			return fmt.Sprintf("failed to marshal params: %v", err), nil
		}
		body = strings.NewReader(string(data))
	}

	req, err := http.NewRequest(t.info.HttpMethod, reqUrl, body)
	if err != nil {
		return fmt.Sprintf("failed to create request: %v", err), nil
	}

	for _, header := range t.info.Headers {
		if header.Key == "" {
			continue
		}

		req.Header.Set(header.Key, header.Value)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("failed to make request: %v", err), nil
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("failed to read response body: %v", err), nil
	}

	return string(respBody), nil
}

func GetAgentTools(agent *repository.Agent) []langchaintools.Tool {
	tools := []langchaintools.Tool{}

	for _, tool := range agent.HttpTools {
		tools = append(tools, &agentTool{info: tool})
	}

	return tools
}
