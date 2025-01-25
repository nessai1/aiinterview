package service

import (
	"fmt"
	"github.com/nessai1/aiinterview/internal/storage"
)

type Service struct {
	config  Config
	storage storage.Storage
}

func NewService(config Config) (*Service, error) {
	s, err := storage.NewPSQLStorageFromAddr(config.PSQLAddress)
	if err != nil {
		return nil, fmt.Errorf("cannot create PSQL storage from address '%s': %w", config.Address, err)
	}

	return &Service{config: config, storage: s}, nil
}

func (s *Service) ListenAndServe() error {

	return nil
}
