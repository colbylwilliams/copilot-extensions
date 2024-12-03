package agent

import (
	"encoding/json"
	"fmt"
)

type ReferenceType string

const (
	ReferenceTypeGitHubRedacted   ReferenceType = "github.redacted"
	ReferenceTypeGitHubAgent      ReferenceType = "github.agent"
	ReferenceTypeGitHubCurrentUrl ReferenceType = "github.current-url"
	ReferenceTypeGitHubFile       ReferenceType = "github.file"
	ReferenceTypeGitHubRepository ReferenceType = "github.repository"
	ReferenceTypeGitHubSnippet    ReferenceType = "github.snippet"
	ReferenceTypeClientFile       ReferenceType = "client.file"
	ReferenceTypeClientSelection  ReferenceType = "client.selection"
)

// Reference copilot reference in chat messages
type Reference struct {
	Type       ReferenceType     `json:"type"`
	ID         string            `json:"id"`
	IsImplicit bool              `json:"is_implicit"`
	Metadata   ReferenceMetadata `json:"metadata"`
	Data       ReferenceData     `json:"-"`
	RawData    json.RawMessage   `json:"data"`
}

// ReferenceMetadata contains metadata about a reference
type ReferenceMetadata struct {
	DisplayName string `json:"display_name"`
	DisplayIcon string `json:"display_icon"`
	DisplayURL  string `json:"display_url"`
}

type ReferenceData interface{}

type ReferenceDataType string

const (
	ReferenceDataTypeGitHubRedacted   ReferenceDataType = "github.redacted"
	ReferenceDataTypeGitHubAgent      ReferenceDataType = "github.agent"
	ReferenceDataTypeGitHubCurrentUrl ReferenceDataType = "github.current-url"
	ReferenceDataTypeGitHubFile       ReferenceDataType = "file"
	ReferenceDataTypeGitHubRepository ReferenceDataType = "repository"
	ReferenceDataTypeGitHubSnippet    ReferenceDataType = "snippet"
	ReferenceDataTypeClientFile       ReferenceDataType = "client.file"
	ReferenceDataTypeClientSelection  ReferenceDataType = "client.selection"
)

// ReferenceDataGitHubRedacted represents a
// redacted copilot reference
type ReferenceDataGitHubRedacted struct {
	Type ReferenceDataType `json:"type"` // "github.redacted"
}

// ReferenceDataGitHubAgent represents a
// reference to a copilot agent
type ReferenceDataGitHubAgent struct {
	AvatarURL string            `json:"avatarURL"`
	ID        int64             `json:"id"`
	Login     string            `json:"login"`
	Type      ReferenceDataType `json:"type"` // "github.agent"
	URL       string            `json:"url"`
}

// ReferenceDataGitHubCurrentUrl represents the
// current URL reference data from a _session message
type ReferenceDataGitHubCurrentUrl struct {
	Type  ReferenceDataType `json:"type"` // "github.current-url" (hardcoded)
	URL   string            `json:"url"`
	Owner string            `json:"owner"`
	Repo  string            `json:"repo"`
	Path  string            `json:"path,omitempty"`
	Hash  string            `json:"hash,omitempty"`
}

// ReferenceDataGitHubFile represents a reference to a
// file on github.com
type ReferenceDataGitHubFile struct {
	CommitOID    string            `json:"commitOID"`
	LanguageID   int64             `json:"languageID"`
	LanguageName string            `json:"languageName"`
	Path         string            `json:"path"`
	Ref          string            `json:"ref"`
	RepoID       int64             `json:"repoID"`
	RepoName     string            `json:"repoName"`
	RepoOwner    string            `json:"repoOwner"`
	Type         ReferenceDataType `json:"type"` // "file"
	URL          string            `json:"url"`
}

// ReferenceDataGitHubRepository represents a
// reference to a repository on github.com
type ReferenceDataGitHubRepository struct {
	CommitOID   string                                   `json:"commitOID"`
	Description string                                   `json:"description"`
	ID          int64                                    `json:"id"`
	Languages   []*ReferenceDataGitHubRepositoryLanguage `json:"languages"`
	Name        string                                   `json:"name"`
	OwnerLogin  string                                   `json:"ownerLogin"`
	OwnerType   string                                   `json:"ownerType"`
	ReadmePath  string                                   `json:"readmePath"`
	Ref         string                                   `json:"ref"`
	RefInfo     ReferenceDataGitHubRepositoryRefInfo     `json:"refInfo"`
	Type        ReferenceDataType                        `json:"type"` // "repository"
	Visibility  string                                   `json:"visibility"`
}
type ReferenceDataGitHubRepositoryLanguage struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}
type ReferenceDataGitHubRepositoryRefInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ReferenceDataGitHubSnippet represents a
// reference to a snippet on github.com
type ReferenceDataGitHubSnippet struct {
	CommitOID    string                          `json:"commitOID"`
	LanguageID   int64                           `json:"languageID"`
	LanguageName string                          `json:"languageName"`
	Path         string                          `json:"path"`
	Range        ReferenceDataGitHubSnippetRange `json:"range"`
	Ref          string                          `json:"ref"`
	RepoID       int64                           `json:"repoID"`
	RepoName     string                          `json:"repoName"`
	RepoOwner    string                          `json:"repoOwner"`
	Type         ReferenceDataType               `json:"type"` // "snippet"
	URL          string                          `json:"url"`
}
type ReferenceDataGitHubSnippetRange struct {
	End   int32 `json:"end"`
	Start int32 `json:"start"`
}

// ReferenceDataClientFile represents a
// reference to a file on a client (e.g. vscode)
type ReferenceDataClientFile struct {
	Content  string            `json:"content"`
	Language string            `json:"language"`
	Type     ReferenceDataType `json:"type"` // "client.file" (hardcoded)
}

// ReferenceDataClientSelection represents a
// reference to a selection on a client (e.g. vscode)
type ReferenceDataClientSelection struct {
	Content string                               `json:"content"`
	End     ReferenceDataClientSelectionLocation `json:"end"`
	Start   ReferenceDataClientSelectionLocation `json:"start"`
	Type    ReferenceDataType                    `json:"type"` // "client.selection" (hardcoded)
}
type ReferenceDataClientSelectionLocation struct {
	Col  int32 `json:"col"`
	Line int32 `json:"line"`
}

func (r *Reference) UnmarshalJSON(data []byte) error {
	type reference Reference

	if err := json.Unmarshal(data, (*reference)(r)); err != nil {
		return err
	}

	var d ReferenceData

	switch r.Type {
	case ReferenceTypeGitHubRedacted:
		d = &ReferenceDataGitHubRedacted{}
	case ReferenceTypeGitHubAgent:
		d = &ReferenceDataGitHubAgent{}
	case ReferenceTypeGitHubCurrentUrl:
		d = &ReferenceDataGitHubCurrentUrl{Type: ReferenceDataTypeGitHubCurrentUrl}
	case ReferenceTypeGitHubFile:
		d = &ReferenceDataGitHubFile{}
	case ReferenceTypeGitHubRepository:
		d = &ReferenceDataGitHubRepository{}
	case ReferenceTypeGitHubSnippet:
		d = &ReferenceDataGitHubSnippet{}
	case ReferenceTypeClientFile:
		d = &ReferenceDataClientFile{Type: ReferenceDataTypeClientFile}
	case ReferenceTypeClientSelection:
		d = &ReferenceDataClientSelection{Type: ReferenceDataTypeClientSelection}
	default:
		return fmt.Errorf("unknown reference type: %s", r.Type)
		// d = &struct{}{}
	}

	if err := json.Unmarshal(r.RawData, d); err != nil {
		return err
	}

	r.Data = d

	return nil
}

func (u *ReferenceDataGitHubCurrentUrl) UnmarshalJSON(data []byte) error {
	type referenceDataGitHubCurrentUrl ReferenceDataGitHubCurrentUrl

	if err := json.Unmarshal(data, (*referenceDataGitHubCurrentUrl)(u)); err != nil {
		return err
	}

	if matches := RepoRe.FindStringSubmatch(u.URL); matches != nil {
		if o := RepoRe.SubexpIndex("owner"); o > -1 {
			u.Owner = matches[o]
		}
		if r := RepoRe.SubexpIndex("repo"); r > -1 {
			u.Repo = matches[r]
		}
		if p := RepoRe.SubexpIndex("path"); p > -1 {
			u.Path = matches[p]
		}
		if h := RepoRe.SubexpIndex("hash"); h > -1 {
			u.Hash = matches[h]
		}
	}

	return nil
}
