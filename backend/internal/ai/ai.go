package ai

import (
	"context"
	"errors"
	"fmt"
	"github.com/nessai1/aiinterview/internal/domain"
	"github.com/nessai1/aiinterview/internal/prompt"
	"github.com/nessai1/aiinterview/internal/storage"
	"github.com/nessai1/aiinterview/internal/utils"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// TODO: DI required model ID
const modelID = openai.GPT4o

const assistantID = "AI_INTERVIEW"

const promptAssistantIntroduction = "assistant_introduction"

type command string

const cmdSkip = command("SKIP")
const cmdNext = command("NEXT")
const cmdChange = command("CHANGE")
const cmdFeedback = command("FEEDBACK")

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
	if err != nil && !errors.Is(storage.ErrNotFound, err) {
		return nil, fmt.Errorf("cannot retrieve assistant from storage: %w", err)
	} else if errors.Is(storage.ErrNotFound, err) {
		assistant = domain.Assistant{ID: assistantID, Model: modelID, ExternalID: ""}
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

	return &Service{config: config, client: client, assistant: assistant, externalAssistant: externalAssistant, promptStorage: promptStorage}, nil
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

var ErrEmptySections = fmt.Errorf("empty sections")
var ErrCorrupt = fmt.Errorf("corrupt message")
var InvalidCommand = fmt.Errorf("invalid command")

// Start Create new thread with context of selected topics and time restriction
func (s *Service) Start(ctx context.Context, topics []domain.Topic, sectionTimeMinutes int) (domain.ChatThread, string, error) {
	secret, err := utils.RandomStringFromCharset(32)
	if err != nil {
		return domain.ChatThread{}, "", fmt.Errorf("cannot generate secret for thread")
	}

	resp, err := s.client.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		return domain.ChatThread{}, "", fmt.Errorf("cannot create thread: %w", err)
	}

	thread := domain.ChatThread{
		ID:     resp.ID,
		Secret: secret,
	}

	topicsStr := ""
	for _, topic := range topics {
		topicsStr += topic.Name + " - " + string(topic.Grade) + "; "
	}

	message, err := s.promptStorage.LoadPrompt("start_cmd", map[string]string{
		"TOPICS":       topicsStr,
		"SECRET":       thread.Secret,
		"SECTION_TIME": strconv.Itoa(sectionTimeMinutes),
	})

	if err != nil {
		return domain.ChatThread{}, "", fmt.Errorf("cannot load start prompt: %w", err)
	}

	respMessage, err := s.send(ctx, thread, message)
	if err != nil {
		return domain.ChatThread{}, "", fmt.Errorf("cannot send start message: %w", err)
	}

	return thread, respMessage, nil
}

func (s *Service) send(ctx context.Context, thread domain.ChatThread, message string) (string, error) {
	_, err := s.client.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	if err != nil {
		return "", fmt.Errorf("cannot create message: %w", err)
	}

	run, err := s.client.CreateRun(ctx, thread.ID, openai.RunRequest{AssistantID: s.assistant.ExternalID})
	if err != nil {
		return "", fmt.Errorf("cannot create run: %w", err)
	}

	for {
		<-time.After(3 * time.Second)
		run, err = s.client.RetrieveRun(ctx, thread.ID, run.ID)
		if err != nil {
			return "", fmt.Errorf("cannot retrieve run: %w", err)
		}

		if run.Status == openai.RunStatusCompleted || run.Status == openai.RunStatusFailed {
			break
		}
	}

	if run.Status != openai.RunStatusCompleted {
		return "", fmt.Errorf("run status is not completed: %s", run.Status)
	}

	l, err := s.client.ListMessage(ctx, thread.ID, nil, nil, nil, nil, nil)

	if err != nil {
		return "", fmt.Errorf("cannot list messages: %w", err)
	}

	respMessage := l.Messages[0].Content[0].Text.Value

	return respMessage, parseErrorMessage(respMessage)
}

func parseErrorMessage(msg string) error {
	switch msg {
	case "INVALID_COMMAND":
		return InvalidCommand
	case "EMPTY_SECTIONS":
		return ErrEmptySections
	case "CORRUPT":
		return ErrCorrupt
	}

	return nil
}

func (s *Service) sendSimpleCommand(ctx context.Context, cmd command, thread domain.ChatThread) (string, error) {
	commandText, err := s.promptStorage.LoadPrompt("simple_cmd", map[string]string{
		"SECRET":  thread.Secret,
		"COMMAND": string(cmd),
	})

	if err != nil {
		return "", fmt.Errorf("cannot load next command prompt: %w", err)
	}

	resp, err := s.send(ctx, thread, commandText)

	if err != nil {
		return "", fmt.Errorf("cannot send '%s' command: %w", cmd, err)
	}

	return resp, nil
}

// Answer on current question in thread. Simple send message in thread
func (s *Service) Answer(ctx context.Context, thread domain.ChatThread, answer string) (string, error) {
	resp, err := s.send(ctx, thread, answer)
	if err != nil {
		return "", fmt.Errorf("cannot send answer: %w", err)
	}

	return resp, nil
}

// Next Get new answer in current section. Return next question text
func (s *Service) Next(ctx context.Context, thread domain.ChatThread) (string, error) {
	return s.sendSimpleCommand(ctx, cmdNext, thread)
}

func (s *Service) Skip(ctx context.Context, thread domain.ChatThread) (string, error) {
	return s.sendSimpleCommand(ctx, cmdSkip, thread)
}

func (s *Service) Change(ctx context.Context, thread domain.ChatThread) (string, error) {
	return s.sendSimpleCommand(ctx, cmdChange, thread)
}

func (s *Service) Feedback(ctx context.Context, thread domain.ChatThread) (string, error) {
	return s.sendSimpleCommand(ctx, cmdFeedback, thread)
}
