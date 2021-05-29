package main

import (
	"context"
	"net/http"
	"search_songs_21/pkg/searcher"

	"github.com/labstack/echo/v4"
)

func main() {
	var s searcher.Searcher
	var de searcher.DocumentExtractor
	ds := searcher.DirSearcher{
		Dir:      "data",
		JobCount: 10,
	}
	if err := ds.Init(); err != nil {
		panic(err)
	}
	s = &ds
	de = &ds
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})
	e.GET("/search", func(c echo.Context) error {
		query := c.QueryParam("q")
		matches := s.Search(context.TODO(), query)
		return c.JSON(http.StatusOK, matches)
		//return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/doc/:id", func(c echo.Context) error {
		docID := c.Param("id")
		res, err := de.GetDocument(context.TODO(), docID)
		if err != nil {
			e.Logger.Errorf("Failed to find doc %s with error %v", docID, err)
			return c.String(http.StatusNotFound, "Not found :(")
		}
		return c.String(http.StatusOK, res)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
