module github.com/colbylwilliams/copilot-extensions

go 1.23.1

require (
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.8.0
	github.com/colbylwilliams/copilot-go v1.0.10
	github.com/go-chi/chi/v5 v5.1.0
	github.com/google/go-github/v67 v67.0.0
	github.com/google/uuid v1.6.0
	github.com/openai/openai-go v0.1.0-alpha.39
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.16.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	golang.org/x/crypto v0.30.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

replace github.com/colbylwilliams/copilot-extensions => ./

// replace github.com/colbylwilliams/copilot-go => ./../copilot-go
