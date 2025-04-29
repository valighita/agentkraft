package sql

import (
	"errors"
	"valighita/agentkraft/repository"

	"gorm.io/gorm"
)

type SQLAgentsRepository struct {
	db *gorm.DB
}

func NewSQLAgentsRepository(db *gorm.DB) repository.AgentsRepository {
	return &SQLAgentsRepository{db: db}
}

func (r *SQLAgentsRepository) GetAgentByID(id uint) (*repository.Agent, error) {
	var agent repository.Agent
	if err := r.db.First(&agent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &agent, nil
}

func (r *SQLAgentsRepository) GetAllAgents() ([]repository.Agent, error) {
	var agents []repository.Agent
	if err := r.db.Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}
