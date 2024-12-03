package dev

import (
	"context"
	"fmt"
	"os"

	octokit "github.com/octokit/go-sdk/pkg"
)

func main() {
	gh, err := octokit.NewApiClient(
		octokit.WithUserAgent("my-thing"),
		octokit.WithTokenAuthentication(os.Getenv("GITHUB_TOKEN")),
	)
	if err != nil {
		fmt.Println("error: ", err)
	}

	ctx := context.Background()

	user, err := gh.User().Get(ctx, nil)
	if err != nil {
		fmt.Println("error: ", err) // panic: index is empty
	}

	fmt.Println("user: ", user)
}
