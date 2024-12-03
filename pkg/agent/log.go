// Package main implements an HTTP service and sets up a server.
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

func GetTempLogFileForRequest(ctx context.Context) (*os.File, error) {
	requestID := middleware.GetReqID(ctx)

	file := fmt.Sprintf("tmp/%s-res.json", requestID)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func WriteRequestToFile(ctx context.Context, r Request) error {
	requestID := middleware.GetReqID(ctx)

	file := fmt.Sprintf("tmp/%s-req.json", requestID)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(r); err != nil {
		return err
	}

	return nil
}

func printJson(o any, newLines bool) {
	if newLines {
		fmt.Println("")
	}

	if j, err := json.MarshalIndent(o, "", "  "); err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println(string(j))
	}

	if newLines {
		fmt.Println("")
	}
}

func PrintJson(o any) {
	printJson(o, false)
}

func PrintJsonWithNewLines(o any) {
	printJson(o, true)
}

func (r *RepoItemRef) Print() {
	PrintJsonWithNewLines(r)
}

func (r *Reference) Print() {
	fmt.Println("")
	switch d := r.Data.(type) {
	case *ReferenceDataGitHubRedacted:
		fmt.Println("data type: ", d.Type)
		PrintJson(r)
	case *ReferenceDataGitHubAgent:
		fmt.Println("data type: ", d.Type)
		PrintJson(r)
	case *ReferenceDataGitHubCurrentUrl:
		fmt.Println("data type: ", d.Type)
		PrintJson(r)
	case *ReferenceDataGitHubFile:
		fmt.Println("data type: ", d.Type)
		PrintJson(r)
	case *ReferenceDataGitHubRepository:
		fmt.Println("data type: ", d.Type)
		PrintJson(r)
	case *ReferenceDataGitHubSnippet:
		fmt.Println("data type: ", d.Type)
		PrintJson(r)
	case *ReferenceDataClientFile:
		fmt.Println("data type: ", d.Type)
		PrintJson(r)
	case *ReferenceDataClientSelection:
		fmt.Println("data type: ", d.Type)
		PrintJson(r)
	default:
		fmt.Println("unknown data type: ", d)
		PrintJson(r)
	}
	fmt.Println("")
}
