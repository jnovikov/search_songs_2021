package searcher

import (
	"context"
)

type Match struct {
	Line         string `json:"line"`
	LineNum      int    `json:"lineNum"`
	DocumentName string `json:"documentName"`
}

type Searcher interface {
	Search(ctx context.Context, query string) []Match
}

type DocumentExtractor interface {
	GetDocument(ctx context.Context, id string) (string, error)
}
