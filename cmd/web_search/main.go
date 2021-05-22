package main

import (
	"context"
	"fmt"
	"net/http"
	"search_songs_21/pkg/searcher"

	"github.com/labstack/echo/v4"
)

func main() {
	var s searcher.Searcher
	ds := searcher.DirSearcher{
		Dir:      "data",
		JobCount: 10,
	}
	if err := ds.Init(); err != nil {
		panic(err)
	}
	s = &ds
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		query := c.QueryParam("q")
		matches := s.Search(context.TODO(), query)
		fmt.Println(matches)
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
