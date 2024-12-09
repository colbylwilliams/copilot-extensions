// Package main implements an HTTP service and sets up a server.
package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/colbylwilliams/copilot-extensions/pkg/agent/agentplateng"
	"github.com/colbylwilliams/copilot-extensions/pkg/github"
	"github.com/colbylwilliams/copilot-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

const (
	envFile      = ".env.plat-eng"
	defaultPort  = "3333"
	readTimeout  = 5 * time.Second   // 5 seconds
	writeTimeout = 300 * time.Second // 5 minutes
)

func main() {
	if err := realMain(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("failed to run service:", err)
		os.Exit(1)
	}
}

//nolint:maintidx // main can have a lot of code.
func realMain() error {
	fmt.Println("Starting api")
	fmt.Println("loading config from", envFile)

	cfg, err := copilot.LoadConfig(envFile)
	if err != nil {
		return err
	}
	if cfg.HTTPPort == "" {
		fmt.Println("no PORT environment variable specified, defaulting to", defaultPort)
		cfg.HTTPPort = defaultPort
	}

	verifier, err := copilot.NewPayloadVerifier()
	if err != nil {
		return fmt.Errorf("failed to create payload authenticator: %w", err)
	}

	// create the azure credential
	azureCredential, err := azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{TenantID: cfg.Azure.TenantID})
	if err != nil {
		return err
	}

	// create the openai client
	oai := openai.NewClient(
		azure.WithEndpoint(cfg.Azure.OpenAIEndpoint, cfg.Azure.OpenAIAPIVersion),
		azure.WithTokenCredential(azureCredential),
	)

	// create the router
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))

	platengAgent := agentplateng.New(cfg, oai)
	platengRoute := platengAgent.Route()

	router.Route("/"+platengRoute, func(r chi.Router) {
		router.Post("/events", github.WebhookHandler)
		router.Post("/agent", copilot.AgentHandler(verifier, platengAgent))
	})

	// TODO: Add auth routes (per agent)
	// authHandlers := &auth.AuthHandlers{
	// 	ClientID: cfg.GitHubAppClientID,
	// 	Callback: cfg.GitHubAppFQDN + "/auth/callback",
	// }

	// router.Route("/auth", func(r chi.Router) {
	// 	r.Get("/authorization", authHandlers.PreAuth)
	// 	r.Get("/callback", authHandlers.PostAuth)
	// })

	addr := ":" + cfg.HTTPPort
	if cfg.IsDevelopment() {
		addr = "127.0.0.1" + addr
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	fmt.Println("Starting server on port " + cfg.HTTPPort)

	return server.ListenAndServe()
}
