package pagesaver

type PageSaver interface {
	Save(file, content string) error
}
