package searcher

import (
	"context"
)


type Match struct {
	Line string
	LineNum int
	DocumentName string
}

type Searcher interface {
	Search(ctx context.Context, query string) []Match
}
