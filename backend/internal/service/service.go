package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nessai1/aiinterview/internal/storage"
	"github.com/nessai1/aiinterview/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

type Service struct {
	config  Config
	storage storage.Storage
	logger  *zap.Logger

	authService *AuthService
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

	if config.InvitationCode == "" {
		config.InvitationCode, err = utils.RandomStringFromCharset(5)
		if err != nil {
			return nil, fmt.Errorf("cannot generate random string for invitation code: %w", err)
		}
	}

	authService := AuthService{secret: config.Secret}

	return &Service{config: config, storage: s, logger: logger, authService: &authService}, nil
}

func (s *Service) ListenAndServe() error {
	s.logger.Info("Service started", zap.Bool("dev", s.config.IsDev), zap.String("address", s.config.Address), zap.String("invitation_code", s.config.InvitationCode))

	err := http.ListenAndServe(s.config.Address, s.buildRouter())
	if err != nil {
		return fmt.Errorf("error while listening http: %w", err)
	}

	return nil
}

func (s *Service) buildRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/join/{invitation}", s.handlePublicAuthorize).Methods("GET")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(s.middlewareTokenAuth)

	apiRouter.HandleFunc("/interviews", s.handleAPIGetInterviewList)

	publicRouter := router.PathPrefix("/").Subrouter()
	publicRouter.Use(s.middlewareTokenAuth)
	publicRouter.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})

	return router
}
