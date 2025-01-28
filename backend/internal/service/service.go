package service

import (
	"fmt"
	"github.com/nessai1/aiinterview/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

type Service struct {
	config  Config
	storage storage.Storage
	logger  *zap.Logger
}

func NewService(config Config) (*Service, error) {
	s, err := storage.NewPSQLStorageFromAddr(config.PSQLAddress)
	if err != nil {
		return nil, fmt.Errorf("cannot create PSQL storage from address '%s': %w", config.Address, err)
	}

	var logger *zap.Logger
	if config.IsDev {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("cannot create development logger: %w", err)
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("cannot create production logger: %w", err)
		}
	}

	return &Service{config: config, storage: s, logger: logger}, nil
}

func (s *Service) ListenAndServe() error {
	s.logger.Info("Service started", zap.Bool("dev", s.config.IsDev), zap.String("address", s.config.Address))
	return nil
}

func (s *Service) buildMux() *http.ServeMux {
	mux := http.NewServeMux()

	// public section
	mux.HandleFunc("/authorize/{token}", s.handlePublicAuthorize)

	// api section
	mux.HandleFunc("/api/interview/list", s.handleAPIGetInterviewList)
}
