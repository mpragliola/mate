package main

import (
	"fmt"
	"path/filepath"

	"github.com/mpragliola/mate/internal/aggregator"
	"github.com/mpragliola/mate/internal/file"
	"github.com/mpragliola/mate/internal/mate"
	"github.com/mpragliola/mate/internal/postsaver"
	"github.com/mpragliola/stopwatch"
)

func main() {
	p := mate.NewProject("example", "posts", "layouts", "public", "tags")
	pa := mate.NewParser()
	ps := &postsaver.FilePostSaver{}
	w := mate.NewWriter(ps)

	stopwatch := stopwatch.NewStart()

	paths, err := aggregator.AggregatePostPaths(p.GetPostsDirectory())
	if err != nil {
		panic(err)
	}

	stopwatch.Mark("Aggregate .md")

	pages, err := pa.ParsePaths(paths, p)
	if err != nil {
		panic(err)
	}

	stopwatch.Mark("Parse to HTML")

	err = w.Clean(p)
	if err != nil {
		panic(err)
	}

	stopwatch.Mark("Wipe public folder")

	err = w.WritePages(pages, p)
	if err != nil {
		panic(err)
	}

	stopwatch.Mark("Write HTML - posts")

	err = w.WriteTags(pages, p)
	if err != nil {
		panic(err)
	}

	stopwatch.Mark("Write HTML - tags")

	err = file.Copy(
		filepath.Join(p.GetLayoutsDirectory(), "assets"),
		filepath.Join(p.GetPublicDirectory(), "assets"),
		ps,
	)
	if err != nil {
		panic(err)
	}

	stopwatch.Mark("Copy assets")

	fmt.Printf("No. of posts:%11d posts\n", len(paths))

	stopwatch.Dump()
}
