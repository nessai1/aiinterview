package ai

import (
	"context"
	"fmt"
	"github.com/nessai1/aiinterview/internal/domain"
	"github.com/nessai1/aiinterview/internal/prompt"
	"github.com/nessai1/aiinterview/internal/storage"
	"github.com/nessai1/aiinterview/internal/utils"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

// TODO: DI required model ID
const modelID = openai.GPT4o

const assistantID = "AI_INTERVIEW"

const promptAssistantIntroduction = "assistant_introduction"

type Config struct {
	Token string

	ProxyURL      string
	ProxyLogin    string
	ProxyPassword string
}

type Service struct {
	client        *openai.Client
	promptStorage *prompt.Storage
	config        Config

	externalAssistant openai.Assistant
	assistant         domain.Assistant
}

func NewService(promptStorage *prompt.Storage, logger *zap.Logger, st storage.Storage, config Config) (*Service, error) {
	// У меня не так много времени, я хз как сделать так,
	// что-бы не было зависимости от storage и лаконично впихнуть поиск/создание нового ассистента
	// поэтому пока DI хранилища :)

	clientConfig := openai.DefaultConfig(config.Token)
	if config.ProxyURL != "" {
		clientConfig.HTTPClient = createProxyClient(config)
	}

	client := openai.NewClientWithConfig(clientConfig)

	// Check if the client is working and check available models
	resp, err := client.ListModels(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("cannot list models: %w", err)
	}

	found := false
	for _, model := range resp.Models {
		if model.ID == modelID {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("required model '%s' not found", modelID)
	}

	assistant, err := st.GetAssistant(context.TODO(), assistantID)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve assistant from storage: %w", err)
	}

	instructions, err := promptStorage.LoadPrompt(promptAssistantIntroduction, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot load assistant introduction prompt: %w", err)
	}

	externalAssistant, err := client.RetrieveAssistant(context.TODO(), assistant.ExternalID)
	if err != nil {
		logger.Error("cannot retrieve assistant from OpenAI, try to create new", zap.Error(err), zap.String("assistant_external_id", assistant.ExternalID))

		externalAssistant, err = client.CreateAssistant(context.TODO(), openai.AssistantRequest{
			Model:        modelID,
			Name:         utils.StringPtr(assistantID),
			Instructions: &instructions,
		})
		if err != nil {
			return nil, fmt.Errorf("cannot create assistant: %w", err)
		}

		assistant.ExternalID = externalAssistant.ID
		err = st.SetAssistant(context.TODO(), assistant)
		if err != nil {
			return nil, fmt.Errorf("cannot save assistant to storage: %w", err)
		}
	}

	return &Service{config: config, client: client, assistant: assistant, externalAssistant: externalAssistant}, nil
}

func createProxyClient(config Config) *http.Client {
	proxyURL, err := url.Parse(config.ProxyURL)
	if err != nil {
		panic(err)
	}

	if config.ProxyLogin != "" && config.ProxyPassword != "" {
		proxyURL.User = url.UserPassword(config.ProxyLogin, config.ProxyPassword)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	return &http.Client{Transport: transport}
}
