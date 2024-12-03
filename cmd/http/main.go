// Package main implements an HTTP service and sets up a server.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/colbylwilliams/copilot-extensions/pkg/agent"
	"github.com/colbylwilliams/copilot-extensions/pkg/agent/agentplateng"
	"github.com/colbylwilliams/copilot-extensions/pkg/agent/payload"
	"github.com/colbylwilliams/copilot-extensions/pkg/auth"
	"github.com/colbylwilliams/copilot-extensions/pkg/config"
	"github.com/colbylwilliams/copilot-extensions/pkg/github"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

const (
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

	const envFile = ".env.plat-eng"
	fmt.Println("loading config from", envFile)

	// load the config
	cfg, err := config.Load(envFile)
	if err != nil {
		return err
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(cfg.HTTPPort)
		fmt.Println("no port specified, will use port from config file:", cfg.HTTPPort)
	}
	fmt.Println("using port:", port)

	// get the public key used to verify the payload signature
	pubKey, err := github.FetchPublicKey()
	if err != nil {
		return fmt.Errorf("failed to fetch public key: %w", err)
	}

	payloadAuthenticator, err := payload.NewAuthenticator(pubKey)
	if err != nil {
		return fmt.Errorf("failed to create payload authenticator: %w", err)
	}

	// create the azure credential
	azureCredential, err := azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{TenantID: cfg.AzureTenantID})
	if err != nil {
		return err
	}

	// create the openai client
	oai := openai.NewClient(
		azure.WithEndpoint(cfg.AzureOpenAIEndpoint, cfg.AzureOpenAIAPIVersion),
		azure.WithTokenCredential(azureCredential),
	)

	platengAgent := agentplateng.New(cfg, oai)

	// create the router
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Get("/_ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	router.Post("/events", github.WebhookHandler)

	router.Post("/agent", executeAgent(payloadAuthenticator, platengAgent))

	authHandlers := &auth.AuthHandlers{
		ClientID: cfg.GitHubAppClientID,
		Callback: cfg.GitHubAppFQDN + "/auth/callback",
	}

	router.Route("/auth", func(r chi.Router) {
		r.Get("/authorization", authHandlers.PreAuth)
		r.Get("/callback", authHandlers.PostAuth)
	})

	addr := ":" + port
	if cfg.IsDevelopment() {
		addr = "127.0.0.1" + addr // Prevents MacOS from prompting you about accepting network connections.
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	fmt.Println("Starting server on port " + port)

	return server.ListenAndServe()
}

type executableAgent interface {
	Execute(ctx context.Context, integrationID, token string, req *agent.Request, w http.ResponseWriter) error
}

func executeAgent(pa payload.Authenticator, a executableAgent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to read request body: %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// print all the headers
		// fmt.Println("headers:")
		// for k, v := range r.Header {
		// 	fmt.Printf("  %s: %s\n", k, v)
		// }

		identifier := r.Header.Get("Github-Public-Key-Identifier")
		sig := r.Header.Get("Github-Public-Key-Signature")
		isValid, err := pa.IsValid(r.Context(), b, identifier, sig)
		if err != nil {
			fmt.Printf("failed to validate payload signature: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !isValid {
			http.Error(w, "invalid payload signature", http.StatusUnauthorized)
			return
		}

		var req agent.Request
		if err := json.Unmarshal(b, &req); err != nil {
			fmt.Printf("failed to unmarshal request: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		apiToken := r.Header.Get("X-GitHub-Token")
		integrationID := r.Header.Get("Copilot-Integration-Id")
		if err := a.Execute(r.Context(), integrationID, apiToken, &req, w); err != nil {
			fmt.Printf("failed to execute agent: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
