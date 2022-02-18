package postsaver

type PostSaver interface {
	Save(file, content string) error
}
