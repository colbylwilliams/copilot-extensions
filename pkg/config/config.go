// package config provides a way to configure the app.
// inspired by: github.com/github/copilot-api/cmd/http/config/config.go
package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
)

const (
	// azure
	AzureTenantIDKey         string = "AZURE_TENANT_ID"
	AzureOpenAIAPIKey        string = "AZURE_OPENAI_API_KEY"
	AzureOpenAIEndpointKey   string = "AZURE_OPENAI_ENDPOINT"
	AzureOpenAIAPIVersionKey string = "OPENAI_API_VERSION"
	// github
	GitHubAppIDKey             string = "GITHUB_APP_ID"
	GitHubAppClientIDKey       string = "GITHUB_APP_CLIENT_ID"
	GitHubAppClientSecretKey   string = "GITHUB_APP_CLIENT_SECRET"
	GitHubAppPrivateKeyPathKey string = "GITHUB_APP_PRIVATE_KEY_PATH"
	GitHubAppUserAgentKey      string = "GITHUB_APP_USER_AGENT"
	GitHubAppWebhookSecretKey  string = "GITHUB_APP_WEBHOOK_SECRET"
	GitHubAppFQDNKey           string = "GITHUB_APP_FQDN"
	// github installation
	GitHubAppDefaultInstallationIDKey string = "GITHUB_APP_DEFAULT_INSTALLATION_ID"
	// chat
	OpenAIChatModelKey string = "OPENAI_CHAT_MODEL"
	// assistant
	OpenAIAssistantModelKey            string = "OPENAI_ASSISTANT_MODEL"
	OpenAIAssistantIDKey               string = "OPENAI_ASSISTANT_ID"
	OpenAIAssistantNameKey             string = "OPENAI_ASSISTANT_NAME"
	OpenAIAssistantDescriptionKey      string = "OPENAI_ASSISTANT_DESCRIPTION"
	OpenAIAssistantInstructionsFileKey string = "OPENAI_ASSISTANT_INSTRUCTIONS_FILE"
)

const (
	// azure
	AzureOpenAIAPIVersionDefault string = "2024-07-01-preview" //"2024-06-01"
	// chat
	OpenAIChatModelDefault string = openai.ChatModelGPT4o
	// assistant
	OpenAIAssistantModelDefault            string = openai.ChatModelGPT4o
	OpenAIAssistantNameDefault             string = "Helpful Assistant"
	OpenAIAssistantDescriptionDefault      string = "A helpful assistant."
	OpenAIAssistantInstructionsFileDefault string = "instructions.sample.md"
)

type Config struct {
	Environment string

	HTTPPort int

	// azure
	AzureTenantID         string
	AzureOpenAIEndpoint   string
	AzureOpenAIAPIVersion string

	// github
	GitHubAppID             int64
	GitHubAppClientID       string
	GitHubAppClientSecret   string
	GitHubAppPrivateKeyPath string
	GitHubAppPrivateKey     []byte
	GitHubAppUserAgent      string
	GitHubAppWebhookSecret  string
	GitHubAppFQDN           string

	// github installation
	GitHubAppDefaultInstallationID int64

	// chat
	ChatModel string

	// assistant
	// If an assistant ID is provided, we will use it. Otherwise, we will create a new assistant.
	// If we create a new assistant, it should be deleted when the app is stopped.
	AssistantIDProvided       bool
	AssistantModel            string
	AssistantID               string
	AssistantName             string
	AssistantDescription      string
	AssistantInstructions     string
	AssistantInstructionsFile string
}

func Load(env ...string) (*Config, error) {
	// Load environment variables from .env files.
	// Load doesn't really return an error, so we ignore it.
	_ = godotenv.Load(env...)

	cfg := &Config{}

	cfg.Environment = getEnvOrDefault("ENVIRONMENT", "development")

	cfg.HTTPPort = 3333

	// azure
	cfg.AzureTenantID = getRequiredEnv(AzureTenantIDKey)
	cfg.AzureOpenAIEndpoint = getRequiredEnv(AzureOpenAIEndpointKey)

	cfg.AzureOpenAIAPIVersion = getEnvOrDefault(AzureOpenAIAPIVersionKey, AzureOpenAIAPIVersionDefault)

	// github
	// cfg.GitHubAppID = getRequiredEnv(GitHubAppIDKey)
	cfg.GitHubAppClientID = getRequiredEnv(GitHubAppClientIDKey)
	cfg.GitHubAppClientSecret = getRequiredEnv(GitHubAppClientSecretKey)
	cfg.GitHubAppPrivateKeyPath = getRequiredEnv(GitHubAppPrivateKeyPathKey)

	// TODO - allow for directly setting the private key with GITHUB_APP_PRIVATE_KEY
	// Read key from pem file
	cfg.GitHubAppPrivateKey = getGitHubPrivateKey(cfg.GitHubAppPrivateKeyPath)

	cfg.GitHubAppUserAgent = getRequiredEnv(GitHubAppUserAgentKey)
	cfg.GitHubAppWebhookSecret = getRequiredEnv(GitHubAppWebhookSecretKey)
	cfg.GitHubAppFQDN = getRequiredEnv(GitHubAppFQDNKey)

	// github installation
	id, ok := os.LookupEnv(GitHubAppDefaultInstallationIDKey)
	if ok {
		instID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			panic(err)
		}
		cfg.GitHubAppDefaultInstallationID = instID
	}

	// chat
	cfg.ChatModel = getEnvOrDefault(OpenAIChatModelKey, OpenAIChatModelDefault)

	// assistant
	cfg.AssistantModel = getEnvOrDefault(OpenAIAssistantModelKey, OpenAIAssistantModelDefault)

	assistantID, ok := os.LookupEnv(OpenAIAssistantIDKey)
	if ok && assistantID != "" {
		cfg.AssistantIDProvided = true
		cfg.AssistantID = assistantID
	} else {
		cfg.AssistantIDProvided = false
		cfg.AssistantID = ""

		// these are only used when creating a new assistant
		cfg.AssistantName = getEnvOrDefault(OpenAIAssistantNameKey, OpenAIAssistantNameDefault)
		cfg.AssistantDescription = getEnvOrDefault(OpenAIAssistantDescriptionKey, OpenAIAssistantDescriptionDefault)

		instructionsFile := getEnvOrDefault(OpenAIAssistantInstructionsFileKey, OpenAIAssistantInstructionsFileDefault)

		cfg.AssistantInstructions = getAssistantInstructions(instructionsFile)
		cfg.AssistantInstructionsFile = instructionsFile
	}

	return cfg, nil
}

// IsProduction returns true if the environment is production.
// We consider staging as production as well.
func (cfg *Config) IsProduction() bool {
	return !cfg.IsDevelopment()
}

// IsDevelopment returns true if the environment is development.
func (cfg *Config) IsDevelopment() bool {
	return cfg.Environment == "development"
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(errors.New("Missing required environment variable: " + key))
	}
	return value
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func getAssistantInstructions(instructionsFile string) string {
	// Read instructions from file
	instructions := ""
	if _, err := os.Stat(instructionsFile); err == nil {
		instructionsBytes, err := os.ReadFile(instructionsFile)
		if err != nil {
			panic(err)
		}
		instructions = string(instructionsBytes)
	}
	return instructions
}

func getGitHubPrivateKey(pemFile string) []byte {
	// Read key from pem file
	if _, err := os.Stat(pemFile); err == nil {
		pemBytes, err := os.ReadFile(pemFile)
		if err != nil {
			panic(err)
		}
		return pemBytes
	}
	panic("GitHub App private key not found")
}
