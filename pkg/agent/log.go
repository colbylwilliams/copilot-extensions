package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/colbylwilliams/copilot-go"
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
func WriteRequestToFile(ctx context.Context, r copilot.Request) error {
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

// Print prints the Reference as JSON
// func PrintReference(r *copilot.Reference) {
// 	fmt.Println("")
// 	switch d := r.Data.(type) {
// 	case *copilot.ReferenceDataGitHubRedacted:
// 		fmt.Println("data type: ", d.Type)
// 		PrintJson(r)
// 	case *copilot.ReferenceDataGitHubAgent:
// 		fmt.Println("data type: ", d.Type)
// 		PrintJson(r)
// 	case *copilot.ReferenceDataGitHubCurrentUrl:
// 		fmt.Println("data type: ", d.Type)
// 		PrintJson(r)
// 	case *copilot.ReferenceDataGitHubFile:
// 		fmt.Println("data type: ", d.Type)
// 		PrintJson(r)
// 	case *copilot.ReferenceDataGitHubRepository:
// 		fmt.Println("data type: ", d.Type)
// 		PrintJson(r)
// 	case *copilot.ReferenceDataGitHubSnippet:
// 		fmt.Println("data type: ", d.Type)
// 		PrintJson(r)
// 	case *copilot.ReferenceDataClientFile:
// 		fmt.Println("data type: ", d.Type)
// 		PrintJson(r)
// 	case *copilot.ReferenceDataClientSelection:
// 		fmt.Println("data type: ", d.Type)
// 		PrintJson(r)
// 	default:
// 		fmt.Println("unknown data type: ", d)
// 		PrintJson(r)
// 	}
// 	fmt.Println("")
// }
