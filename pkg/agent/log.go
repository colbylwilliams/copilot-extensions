package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

// GetTempLogFileForRequest returns a temporary file for logging the request
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

// WriteRequestToFile writes the request to a file
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

// PrintJson prints the object as JSON
func PrintJson(o any) {
	printJson(o, false)
}

// PrintJsonWithNewLines prints the object as JSON
// with new lines before and after
func PrintJsonWithNewLines(o any) {
	printJson(o, true)
}

// Print prints the RepoItemRef as JSON
func (r *RepoItemRef) Print() {
	PrintJsonWithNewLines(r)
}

// Print prints the Reference as JSON
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
