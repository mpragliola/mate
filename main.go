package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mpragliola/mate/internal/aggregator"
	"github.com/mpragliola/mate/internal/file"
	"github.com/mpragliola/mate/internal/mate"
	"github.com/mpragliola/mate/internal/postsaver"
	"github.com/mpragliola/mate/internal/server"
	"github.com/mpragliola/stopwatch"
)

func main() {
	p := mate.NewProject("example", "posts", "layouts", "public", "tags")
	pa := mate.NewParser()
	ps := &postsaver.FilePostSaver{}
	w := mate.NewWriter(ps)

	if len(os.Args) > 1 {
		command := os.Args[1]

		switch command {
		case "init":
			if len(os.Args) == 2 {

			}
			return
		case "serve":
			server.Serve(p, pa)
		}
	}

	stopwatch := stopwatch.NewStart()

	paths, err := aggregator.AggregatePostPaths(p)
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

	stopwatch.Mark("Write HTML")

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
