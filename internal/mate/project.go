package mate

import (
	"path/filepath"
	"sort"
)

type Project struct {
	workingDir  string
	postsPath   string
	posts       []Post
	layoutsPath string
	publicPath  string
	tagsPath    string
	tags        map[string]bool
}

func NewProject(workingDir, postsPath, layoutsPath, publicPath, tagsPath string) *Project {
	return &Project{
		workingDir:  workingDir,
		postsPath:   postsPath,
		posts:       make([]Post, 0),
		layoutsPath: layoutsPath,
		publicPath:  publicPath,
		tagsPath:    tagsPath,
		tags:        make(map[string]bool),
	}
}

func (p *Project) AddPost(post Post) {
	p.posts = append(p.posts, post)
}

func (p *Project) Posts() []Post {
	return p.posts
}

func (p *Project) PostsSorted(sortFunction PageSorter) []Post {
	sort.Slice(p.posts, sortFunction)

	return p.posts
}

func (p *Project) AddTags(tag ...string) {
	for _, t := range tag {
		if _, ok := p.tags[t]; !ok {
			p.tags[t] = true
		}
	}
}

func (p *Project) GetPublicTagsPath() string {
	return filepath.Join(
		p.tagsPath,
	)
}

func (p *Project) GetTags() []string {
	tags := make([]string, 0, len(p.tags))
	for tag := range p.tags {
		tags = append(tags, tag)
	}

	return tags
}

func (p *Project) GetPostsDirectory() string {
	return filepath.Join(p.workingDir, p.postsPath)
}

func (p *Project) GetLayoutsDirectory() string {
	return filepath.Join(p.workingDir, p.layoutsPath)
}

func (p *Project) GetPublicDirectory() string {
	return filepath.Join(p.workingDir, p.publicPath)
}

//

func (p *Project) CreatedOnSorter() PageSorter {
	return func(i, j int) bool {
		return p.posts[i].CreatedOn.Before(p.posts[j].CreatedOn)
	}
}

//

type PageSorter func(i, j int) bool
