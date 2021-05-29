package searcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDirSearcher_searchm(t *testing.T) {
	tc := []struct {
		name  string
		input string
		query string
		res   []Match
	}{
		{
			name:  "Simple test",
			input: "Hello\nWorld\nKek\nWorld",
			query: "Kek",
			res: []Match{{
				Line:         "Kek",
				LineNum:      3,
				DocumentName: "",
			}},
		},
		{
			name:  "All match test",
			input: "Hello\nHello\n",
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
			name:  "Empty test",
			input: "",
			query: "query",
			res:   nil,
		},
		{
			name:  "Lowercase test",
			input: "heLlO WorLd\n",
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

func BenchmarkDirSearcher_search(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := DirSearcher{}
		r := strings.NewReader("Hello\nWorld\nLol\nKek\nCheburek\n")
		res := ds.search(r, "lol")
		if len(res) < 1 {
			b.Fatalf("Wrong search result.")
		}
	}
}

func BenchmarkDirSearcher_searchRE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := DirSearcher{}
		r := strings.NewReader("Hello\nWorld\nlol\nKek\nCheburek\n")

		reg, err := regexp.Compile("lol")
		if err != nil {
			b.Fatalf("Failed to compile regexp: %v", err)
		}
		res := ds.searchRE(r, reg)
		if len(res) < 1 {
			b.Fatalf("Wrong search result.")
		}
	}
}

func BenchmarkDirSearcher_Search(b *testing.B) {
	name, err := ioutil.TempDir("", "testdata")
	if err != nil {
		b.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(name)
	match := "kek"
	numFiles := 2000
	numLines := 500
	for i := 0; i < numFiles; i++ {
		f, err := os.Create(path.Join(name, fmt.Sprintf("%d.txt", i)))
		if err != nil {
			b.Fatalf("Failed to create file %s = %v", fmt.Sprintf("%d.txt", i), err)
		}
		for j := 0; j < numLines; j++ {
			f.WriteString("!\n")
		}
		f.WriteString(match)
		if err := f.Close(); err != nil {
			b.Fatalf("Failed to close file %v with error %v", f.Name(), err)
		}
	}

	ds := DirSearcher{
		Dir:      name,
		JobCount: 4,
	}
	if err := ds.Init(); err != nil {
		b.Fatalf("Failed to init searcher with error = %v", err)
	}
	for i := 0; i < b.N; i++ {
		res := ds.Search(context.TODO(), match)
		if len(res) < numFiles {
			b.Fatalf("Wrong search result.")
		}
	}
}

//func Test_dirsearch_search(t *testing.T) {
//
//}
