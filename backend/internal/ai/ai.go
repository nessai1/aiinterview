package ai

import (
	"context"
	"fmt"
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

const assistantName = "AI_INTERVIEW"

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
	assistant     *openai.Assistant
	config        Config
}

func NewService(promptStorage *prompt.Storage, logger *zap.Logger, storage storage.Storage, config Config) (*Service, error) {
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

	assistants, err := client.ListAssistants(context.TODO(), nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot list assistants: %w", err)
	}

	var serviceAssistant *openai.Assistant
	for _, assistant := range assistants.Assistants {
		if assistant.Model == modelID && *assistant.Name == assistantName {
			assistant.ID
			serviceAssistant = &assistant
			break
		}
	}

	if serviceAssistant == nil {
		return nil, fmt.Errorf("required assistant %s not found", assistantName)
	}

	if !found {

		instruction, err := promptStorage.LoadPrompt(promptAssistantIntroduction, map[string]string{})

		if err != nil {
			return nil, fmt.Errorf("cannot load assistant introduction prompt: %w", err)
		}

		client.CreateAssistant(context.TODO(), openai.AssistantRequest{
			Model:        modelID,
			Name:         utils.StringPtr(assistantName),
			Instructions: nil,
		})
	}

	return &Service{config: config, client: client, assistant: serviceAssistant}, nil
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
