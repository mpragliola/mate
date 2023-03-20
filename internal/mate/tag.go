package mate

// Tag ...
type Tag struct {
	Name string
}

// TagList ...
type TagList struct {
	tags map[string]bool
}

func NewTagList() *TagList {
	return &TagList{
		tags: make(map[string]bool),
	}
}

// AddTag ...
func (tl *TagList) AddTag(t string) *TagList {
	tl.tags[t] = true

	return tl
}

// AddTags ...
func (tl *TagList) AddTags(tags []string) *TagList {
	for _, t := range tags {
		tl.AddTag(t)
	}

	return tl
}

// GetTags ...
func (tl *TagList) GetTags() []Tag {
	tags := make([]Tag, 0)
	for v := range tl.tags {
		tags = append(tags, Tag{v})
	}

	return tags
}
