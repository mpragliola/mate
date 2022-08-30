package mate

import (
	"path/filepath"
	"sort"
)

// Project ...
type Project struct {
	workingDir  string
	postsPath   string
	posts       []Post
	layoutsPath string
	publicPath  string
	tagsPath    string
	tags        TagList
}

// NewProject ...
func NewProject(workingDir, postsPath, layoutsPath, publicPath, tagsPath string) *Project {
	return &Project{
		workingDir:  workingDir,
		postsPath:   postsPath,
		posts:       make([]Post, 0),
		layoutsPath: layoutsPath,
		publicPath:  publicPath,
		tagsPath:    tagsPath,
		tags:        TagList{},
	}
}

// AddPost ...
func (p *Project) AddPost(post Post) {
	p.posts = append(p.posts, post)
}

// Posts ...
func (p *Project) Posts() []Post {
	return p.posts
}

// PostsSorted ...
func (p *Project) PostsSorted(sortFunction PageSorter) []Post {
	sort.Slice(p.posts, sortFunction)

	return p.posts
}

// AddTags ...
func (p *Project) AddTags(tag ...string) {
	p.tags.AddTags(tag)
}

// GetPublicTagsPath ...
func (p *Project) GetPublicTagsPath() string {
	return filepath.Join(
		p.GetPublicDirectory(),
		p.tagsPath,
	)
}

// GetTags ...
func (p *Project) GetTags() []Tag {
	return p.tags.GetTags()
}

// GetPostsDirectory ...
func (p *Project) GetPostsDirectory() string {
	return filepath.Join(p.workingDir, p.postsPath)
}

// GetLayoutsDirectory ...
func (p *Project) GetLayoutsDirectory() string {
	return filepath.Join(p.workingDir, p.layoutsPath)
}

// GetPublicDirectory ...
func (p *Project) GetPublicDirectory() string {
	return filepath.Join(p.workingDir, p.publicPath)
}

//

// CreatedOnSorter ...
func (p *Project) CreatedOnSorter() PageSorter {
	return func(i, j int) bool {
		return p.posts[i].CreatedOn.Before(p.posts[j].CreatedOn)
	}
}

//

// PageSorter ...
type PageSorter func(i, j int) bool
