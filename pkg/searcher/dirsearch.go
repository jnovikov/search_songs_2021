package searcher

import (
	"bufio"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

var NegativeJobCountError = errors.New("job count should be positive")

type DirSearcher struct {
	Dir      string
	JobCount int
	fileList []string
	sem      chan struct{}
}

func (ds *DirSearcher) Init() error {
	if ds.JobCount <= 0 {
		return NegativeJobCountError
	}
	fl, err := ds.scanDir()
	if err != nil {
		return err
	}
	ds.fileList = fl
	ds.sem = make(chan struct{}, ds.JobCount)
	return nil
}

func (ds *DirSearcher) search(r io.Reader, q string) (res []Match) {
	scan := bufio.NewScanner(r)
	ln := 1
	for scan.Scan() {
		line := scan.Text()
		if strings.Contains(strings.ToLower(line), strings.ToLower(q)) {
			res = append(res, Match{
				Line:         line,
				LineNum:      ln,
				DocumentName: "",
			})
		}
		ln++
	}
	return res
}

func (ds *DirSearcher) searchRE(r io.Reader, re *regexp.Regexp) (res []Match) {
	scan := bufio.NewScanner(r)
	ln := 1
	for scan.Scan() {
		line := scan.Text()
		if re.Match([]byte(line)) {
			res = append(res, Match{
				Line:         line,
				LineNum:      ln,
				DocumentName: "",
			})
		}
		ln++
	}
	return res
}

func (ds *DirSearcher) scanDir() ([]string, error) {
	infos, err := ioutil.ReadDir(ds.Dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, info := range infos {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
	}
	return files, nil
}

func (ds *DirSearcher) Search(ctx context.Context, query string) (res []Match) {
	var wg sync.WaitGroup
	var m sync.Mutex

Loop:
	for _, fname := range ds.fileList {
		select {
		case ds.sem <- struct{}{}:
			wg.Add(1)
			go func(f string) {
				defer func() {
					wg.Done()
					<-ds.sem
				}()
				r, err := os.Open(path.Join(ds.Dir, f))
				if err != nil {
					log.Printf("Failed to open file %s = %v", f, err)
					return
				}
				searchRes := ds.search(r, query)
				for i := range searchRes {
					searchRes[i].DocumentName = f
				}
				m.Lock()
				res = append(res, searchRes...)
				m.Unlock()
			}(fname)
		case <-ctx.Done():
			break Loop
		}
	}
	wg.Wait()
	return res
}

func (ds *DirSearcher) GetDocument(ctx context.Context, id string) (string, error) {
	// Передать название любого файла и прочитать ваш файлы.
	r, err := os.Open(path.Join(ds.Dir, id))
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
