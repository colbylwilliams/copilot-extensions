package sse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/colbylwilliams/copilot-extensions/pkg/agent"
	"github.com/jclem/sseparser"
)

const (
	sseDataField  = "data"
	sseEventField = "event"
)

type dataEmitter func(data any)

// Parser is a parser for Server-Sent Events (SSE).
type Parser struct {
	buf io.Reader
	fn  dataEmitter
}

// NewParser creates a new SSEParser.
func NewParser(buf io.Reader, fn dataEmitter) *Parser {
	return &Parser{
		buf: buf,
		fn:  fn,
	}
}

// ParseAndEmit parses the SSE stream and emits the parsed events.
func (p *Parser) ParseAndEmit(ctx context.Context) error {
	scanner := sseparser.NewStreamScanner(p.buf)

	for {
		event, _, err := scanner.Next()
		if err != nil {
			if errors.Is(err, sseparser.ErrStreamEOF) {
				return nil
			}
			return fmt.Errorf("failed to read from stream: %w", err)
		}

		eventFields := map[string]string{}
		dataFields := []string{}

		for _, field := range event.Fields() {
			switch field.Name {
			case sseEventField:
				eventFields[field.Name] = field.Value
			case sseDataField:
				dataFields = append(dataFields, field.Value)
			}
		}

		switch {
		case eventFields[sseEventField] == "copilot_confirmation":
			err := emitConfirmation(dataFields, p.fn)
			if err != nil {
				return fmt.Errorf("failed to process confirmation: %w", err)
			}

		case eventFields[sseEventField] == "copilot_references":
			err := emitReferences(dataFields, p.fn)
			if err != nil {
				return fmt.Errorf("failed to process references: %w", err)
			}

		case eventFields[sseEventField] == "copilot_errors":
			err := emitErrors(dataFields, p.fn)
			if err != nil {
				return fmt.Errorf("failed to process errors: %w", err)
			}

		default:
			err = emitDatas(dataFields, p.fn)
			if err != nil {
				return fmt.Errorf("failed to process data: %w", err)
			}
		}
	}
}

func emitErrors(data []string, fn dataEmitter) error {
	for _, d := range data {
		var errs []agent.Error
		if err := json.Unmarshal([]byte(d), &errs); err != nil {
			return fmt.Errorf("failed to unmarshal references: %w", err)
		}

		fn(errs)
	}

	return nil
}

func emitReferences(data []string, fn dataEmitter) error {
	for _, d := range data {
		var refs []agent.Reference
		if err := json.Unmarshal([]byte(d), &refs); err != nil {
			return fmt.Errorf("failed to unmarshal references: %w", err)
		}

		fn(refs)
	}

	return nil
}

func emitConfirmation(data []string, fn dataEmitter) error {
	for _, d := range data {
		var confirmation agent.Confirmation
		if err := json.Unmarshal([]byte(d), &confirmation); err != nil {
			return fmt.Errorf("failed to unmarshal confirmation: %w", err)
		}

		fn(confirmation)
	}

	return nil
}

func emitDatas(datas []string, fn dataEmitter) error {
	for _, data := range datas {
		if data == "" || data == "[DONE]" {
			continue
		}

		var response agent.Response
		if err := json.Unmarshal([]byte(data), &response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		fn(response)
	}

	return nil
}
