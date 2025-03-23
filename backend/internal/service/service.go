package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nessai1/aiinterview/internal/ai"
	"github.com/nessai1/aiinterview/internal/interview"
	"github.com/nessai1/aiinterview/internal/message"
	"github.com/nessai1/aiinterview/internal/prompt"
	"github.com/nessai1/aiinterview/internal/storage"
	"github.com/nessai1/aiinterview/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

type Service struct {
	config Config
	logger *zap.Logger

	authService      *AuthService
	interviewService *interview.Service
	storage          storage.Storage

	messageParser *message.Parser
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

	// TODO: language switch
	promptStorage, err := prompt.NewStorage("ru")
	if err != nil {
		return nil, fmt.Errorf("cannot create prompt storage: %w", err)
	}

	aiService, err := ai.NewService(promptStorage, logger, s, config.OpenAI)
	if err != nil {
		return nil, fmt.Errorf("cannot create AI service: %w", err)
	}

	messageParser := message.NewParser(message.NewHighlighter())

	interviewService, err := interview.NewService(s, aiService, logger, messageParser)
	if err != nil {
		return nil, fmt.Errorf("cannot create interview service: %w", err)
	}

	return &Service{config: config, interviewService: interviewService, storage: s, logger: logger, authService: &authService, messageParser: messageParser}, nil
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

	if s.config.IsDev {
		s.logger.Info("CORS policy disabled for API")
		apiRouter.Use(s.corsAllowMiddleware)
	}

	apiRouter.Use(s.middlewareTokenAuth)

	apiRouter.HandleFunc("/interview/list", s.handleAPIGetInterviewList).Methods("GET")
	apiRouter.HandleFunc("/interview/{interviewID}", s.handleAPIGetInterview).Methods("GET")
	apiRouter.HandleFunc("/interview", s.handleAPICreateInterview).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/preview", s.handleAPIPreviewMessage).Methods("POST", "OPTIONS")

	publicRouter := router.PathPrefix("/").Subrouter()
	publicRouter.Use(s.middlewareTokenAuth)
	publicRouter.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})

	return router
}
