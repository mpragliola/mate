package main

import (
	"fmt"
	"path/filepath"

	"github.com/mpragliola/mate/internal/file"
	"github.com/mpragliola/mate/internal/mate"
	"github.com/mpragliola/mate/internal/pagesaver"
	"github.com/mpragliola/stopwatch"
)

func main() {
	stopwatch := stopwatch.NewStart()

	p := mate.NewProject("example", "posts", "layouts", "public", "tags")

	a := mate.NewAggregator()
	pa := mate.NewParser()
	ps := &pagesaver.FilePageSaver{}
	w := mate.NewWriter(ps)

	paths, err := a.AggregatePostPaths(p)
	if err != nil {
		panic(err)
	}

	stopwatch.Mark("Aggregated .md")

	pages, err := pa.ParsePaths(paths, p)
	if err != nil {
		panic(err)
	}

	stopwatch.Mark("Parsed to HTML")

	err = w.WritePages(pages, p)
	if err != nil {
		panic(err)
	}

	stopwatch.Mark("Written HTML")

	err = file.Copy(
		filepath.Join(
			p.GetLayoutsDirectory(),
			"assets",
		),
		filepath.Join(
			p.GetPublicDirectory(),
			"assets",
		),
		ps,
	)
	if err != nil {
		panic(err)
	}

	stopwatch.Mark("Assets copied")

	fmt.Printf("No. of posts:%11d posts\n", len(paths))

	stopwatch.Dump()
}
