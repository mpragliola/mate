package mate

import (
	"fmt"
	"strings"

	"github.com/gernest/front"
)

// DefaultDelimiter is the default delimiter that encloses the YAML preamble
const DefaultDelimiter = "---"

// FrontMatter will help access the front matter of the posts; the front matter is a YAML
// preamble enclosed by delimiters (`---`) that can provide additional metadata about the
// post itself (f. ex. set anouther layout, provide tags, ...) as key-value pairs.
// It's optional but if present it must be at the beginning of the document.
type FrontMatter struct {
	fm   map[string]string // front matter key-value map
	body string            // actual source body without the FrontMatter
}

func NewFrontMatterFromSource(src string) *FrontMatter {
	m := front.NewMatter()
	m.Handle(DefaultDelimiter, front.YAMLHandler)

	// If a front matter is not parsable for any reason, it just returns
	// an empty front matter and the whole post body
	f, body, err := m.Parse(strings.NewReader(src))
	if err != nil {
		body = string(src)
	}

	return &FrontMatter{
		fm:   convertInterfaceMapToStringMap(f),
		body: body,
	}
}

func (f *FrontMatter) GetValues() map[string]string {
	return f.fm
}

// GetValue will provide the value of a frontmatter field, or a default if not found.
func (f *FrontMatter) GetValue(key, defaultVal string) string {
	if val, ok := f.fm[key]; ok {
		return val
	}

	return defaultVal
}

// GetArrayValue will provide the value of a frontmatter field as a string array, using
// comma as a separator and trimming the whitespace from each string. The default will be
// an empty array.
func (f *FrontMatter) GetArrayValue(key string) []string {
	values := strings.Split(f.GetValue(key, ""), ",")

	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}

	return values
}

// GetBody will provide the contents of the post with front matter removed.
func (f *FrontMatter) GetBody() string {
	return f.body
}

func convertInterfaceMapToStringMap(f map[string]interface{}) map[string]string {
	fm := make(map[string]string, len(f))
	for key, value := range f {
		fm[key] = fmt.Sprintf("%v", value)
	}

	return fm
}
