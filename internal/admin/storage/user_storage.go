package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nanmenkaimak/github-gist/internal/admin/entity"
	"github.com/nanmenkaimak/github-gist/internal/admin/repository"
	"go.uber.org/zap"
)

type DataStorage struct {
	interval time.Duration
	service  repository.User
	users    map[string]*[]entity.User
	mu       sync.RWMutex
	logger   *zap.SugaredLogger
}

func NewDataStorage(interval time.Duration, service repository.User, logger *zap.SugaredLogger) *DataStorage {
	return &DataStorage{
		interval: interval,
		service:  service,
		users:    make(map[string]*[]entity.User),
		mu:       sync.RWMutex{},
		logger:   logger,
	}
}

func (pm *DataStorage) Run() {
	ticker := time.NewTicker(pm.interval)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			startTime := time.Now()

			pm.LoadData()

			elapsedTime := time.Since(startTime)

			timeToNextTick := pm.interval - elapsedTime

			time.Sleep(timeToNextTick)
		}
	}

}

func (pm *DataStorage) LoadData() {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	allUsers, err := pm.service.GetAllUsers()
	if err != nil {
		pm.logger.Errorf("failed to GetAllUsers err: %v", err)
	}

	method := make(map[string]*[]entity.User)

	method["users"] = allUsers

	pm.users = method
}

func (pm *DataStorage) GetAllUsers() (*[]entity.User, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	users, ok := pm.users["users"]
	if !ok {
		return nil, fmt.Errorf("get all users from storage err")
	}
	return users, nil
}
