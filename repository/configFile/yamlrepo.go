package configfile

import "valighita/agentkraft/repository"

func GetYamlRepositories() (repository.AgentsRepository, error) {
	agentRepo, err := NewYAMLAgentsRepository()
	if err != nil {
		return nil, err
	}

	return agentRepo, nil
}
