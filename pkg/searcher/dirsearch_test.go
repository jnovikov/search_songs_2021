package searcher

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDirSearcher_searchm(t *testing.T) {
	tc := []struct {
		name string
		input string
		query string
		res []Match
	} {
		{
			name: "Simple test",
			input:"Hello\nWorld\nKek\nWorld",
			query: "Kek",
			res: []Match{{
				Line:         "Kek",
				LineNum:      3,
				DocumentName: "",
			}},
		},
		{
			name: "All match test",
			input:"Hello\nHello\n",
			query: "Hello",
			res: []Match{
				{
					Line:         "Hello",
					LineNum:      1,
					DocumentName: "",
				},
				{
					Line:         "Hello",
					LineNum:      2,
					DocumentName: "",
				},
			},
		},
		{
			name: "Empty test",
			input:"",
			query: "query",
			res: nil,
		},
		{
			name: "Lowercase test",
			input:"heLlO WorLd\n",
			query: "hello world",
			res: []Match{{
				Line:         "heLlO WorLd",
				LineNum:      1,
				DocumentName: "",
			}},
		},


	}
	for _, c := range tc {
		ds := DirSearcher{}
		r := strings.NewReader(c.input)
		res := ds.search(r, c.query)
		if len(res) != len(c.res) {
			t.Errorf("%s search() missmatch. Wanted docs %d, got %d", c.name, len(c.res), len(res))
		}
		if diff := cmp.Diff(c.res, res); diff != "" {
			t.Errorf("%s search() result mismatch (-want +got):\n%s", c.name, diff)
		}
	}
}

func TestDirSearcher_Search(t *testing.T) {
	ds := DirSearcher{
		Dir:      "testdata",
		JobCount: 5,
	}
	if err := ds.Init(); err != nil {
		t.Fatalf("Failed to init searcher with error = %v", err)
	}
	want := []Match{{
		Line:         "hello",
		LineNum:      1,
		DocumentName: "first.txt",
	}}
	res := ds.Search(context.TODO(), "hello")
	if len(res) != len(want) {
		t.Errorf("Searcn() missmatch. Wanted docs %d, got %d", len(want), len(res))
	}
	if diff := cmp.Diff(want, res); diff != "" {
		t.Errorf("Search() result mismatch (-want +got):\n%s", diff)
	}


}
//func Test_dirsearch_search(t *testing.T) {
//
//}
